package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/volcengine/volcengine-go-sdk/service/ark"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/request"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

const managedKey = "ark_sync_managed"

const (
	modelTypeFoundation = "foundation"
	modelTypeCustom     = "custom"
	modelTypeService    = "service"

	endpointSourceBuiltIn = "builtin"
	endpointSourceUser    = "user"
)

type config struct {
	accessKey, secretKey, region, project, arkAPIKey string
	liteLLMURL, liteLLMMasterKey                     string
	timeout                                          time.Duration
}

type deployment struct {
	ModelName     string         `json:"model_name"`
	LiteLLMParams map[string]any `json:"litellm_params"`
	ModelInfo     map[string]any `json:"model_info"`
}

type action struct {
	kind  string
	model deployment
	id    string
}

type skipStats struct {
	nonRunning, uncallable, shadowedCustom, duplicate int
}

type managedEndpoint struct {
	EndpointModelType *string                                   `type:"string" json:",omitempty"`
	Id                *string                                   `type:"string" json:",omitempty"`
	ModelId           *string                                   `type:"string" json:",omitempty"`
	ModelReference    *ark.ModelReferenceForListEndpointsOutput `type:"structure" json:",omitempty"`
	Name              *string                                   `type:"string" json:",omitempty"`
	ProjectName       *string                                   `type:"string" json:",omitempty"`
	Status            *string                                   `type:"string" json:",omitempty"`
}

type managedEndpointsOutput struct {
	Items      []*managedEndpoint `type:"list" json:",omitempty"`
	PageNumber *int32             `type:"int32" json:",omitempty"`
	PageSize   *int32             `type:"int32" json:",omitempty"`
	TotalCount *int32             `type:"int32" json:",omitempty"`
}

type liteLLMClient struct {
	baseURL, masterKey, arkAPIKey string
	httpClient                    *http.Client
}

func main() {
	apply := flag.Bool("apply", false, "apply changes; default is preview only")
	flag.Parse()

	if err := loadDotEnv(); err != nil {
		exit(err)
	}
	cfg, err := loadConfig()
	if err != nil {
		exit(err)
	}

	ctx := context.Background()
	desired, skipped, err := listARKDeployments(ctx, cfg)
	if err != nil {
		exit(fmt.Errorf("list ARK endpoints: %w", err))
	}
	client := &liteLLMClient{
		baseURL:   strings.TrimRight(cfg.liteLLMURL, "/"),
		masterKey: cfg.liteLLMMasterKey,
		arkAPIKey: cfg.arkAPIKey,
		httpClient: &http.Client{
			Timeout: cfg.timeout,
		},
	}
	if err := syncModels(ctx, client, desired, *apply, os.Stdout); err != nil {
		exit(err)
	}
	if skipped.total() > 0 {
		fmt.Printf(
			"skipped ARK endpoints: non-running=%d, uncallable=%d, shadowed-custom=%d, duplicate=%d\n",
			skipped.nonRunning,
			skipped.uncallable,
			skipped.shadowedCustom,
			skipped.duplicate,
		)
	}
}

func loadDotEnv() error {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	return nil
}

func loadConfig() (config, error) {
	cfg := config{
		accessKey:        os.Getenv("VOLCENGINE_ACCESS_KEY_ID"),
		secretKey:        os.Getenv("VOLCENGINE_SECRET_ACCESS_KEY"),
		region:           valueOr(os.Getenv("VOLCENGINE_REGION"), "cn-beijing"),
		project:          os.Getenv("VOLCENGINE_PROJECT_NAME"),
		arkAPIKey:        os.Getenv("ARK_API_KEY"),
		liteLLMURL:       os.Getenv("LITELLM_BASE_URL"),
		liteLLMMasterKey: os.Getenv("LITELLM_MASTER_KEY"),
		timeout:          30 * time.Second,
	}
	required := map[string]string{
		"VOLCENGINE_ACCESS_KEY_ID":     cfg.accessKey,
		"VOLCENGINE_SECRET_ACCESS_KEY": cfg.secretKey,
		"ARK_API_KEY":                  cfg.arkAPIKey,
		"LITELLM_BASE_URL":             cfg.liteLLMURL,
		"LITELLM_MASTER_KEY":           cfg.liteLLMMasterKey,
	}
	var missing []string
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	if raw := os.Getenv("HTTP_TIMEOUT_SECONDS"); raw != "" {
		seconds, err := strconv.Atoi(raw)
		if err != nil || seconds <= 0 {
			return config{}, fmt.Errorf("HTTP_TIMEOUT_SECONDS must be a positive integer")
		}
		cfg.timeout = time.Duration(seconds) * time.Second
	}
	return cfg, nil
}

func listARKDeployments(ctx context.Context, cfg config) ([]deployment, skipStats, error) {
	httpClient := &http.Client{Timeout: cfg.timeout}
	sess, err := session.NewSession(volcengine.NewConfig().
		WithRegion(cfg.region).
		WithHTTPClient(httpClient).
		WithCredentials(credentials.NewStaticCredentials(cfg.accessKey, cfg.secretKey, "")))
	if err != nil {
		return nil, skipStats{}, err
	}
	client := ark.New(sess)
	deployments, skipped, err := listManagedARKDeployments(ctx, client, cfg)
	if err != nil {
		return nil, skipped, fmt.Errorf("list built-in ARK endpoints: %w", err)
	}

	for page := int32(1); ; page++ {
		input := (&ark.ListEndpointsInput{}).SetPageNumber(page).SetPageSize(100)
		if cfg.project != "" {
			input.SetProjectName(cfg.project)
		}
		out, err := client.ListEndpointsWithContext(ctx, input)
		if err != nil {
			return nil, skipped, err
		}
		for _, endpoint := range out.Items {
			if endpoint == nil || !strings.EqualFold(volcengine.StringValue(endpoint.Status), "Running") {
				skipped.nonRunning++
				continue
			}
			modelType := arkModelType(endpoint)
			modelName := arkModelName(endpoint, modelType)
			inferenceID, inferenceType := arkInferenceTarget(endpoint, modelType)
			if modelName == "" || inferenceID == "" {
				skipped.uncallable++
				continue
			}
			project := valueOr(volcengine.StringValue(endpoint.ProjectName), cfg.project)
			deployments = append(deployments, desiredDeployment(
				modelName,
				volcengine.StringValue(endpoint.Id),
				inferenceID,
				inferenceType,
				project,
				cfg.arkAPIKey,
				modelType,
				endpointSourceUser,
			))
		}
		if len(out.Items) == 0 || out.TotalCount != nil && page*100 >= *out.TotalCount || out.TotalCount == nil && len(out.Items) < 100 {
			break
		}
	}
	deployments, skipped.shadowedCustom, skipped.duplicate = normalizeDeployments(deployments)
	sort.Slice(deployments, func(i, j int) bool {
		return managedID(deployments[i]) < managedID(deployments[j])
	})
	return deployments, skipped, nil
}

func listManagedARKDeployments(ctx context.Context, client *ark.ARK, cfg config) ([]deployment, skipStats, error) {
	var deployments []deployment
	var skipped skipStats
	for page := int32(1); ; page++ {
		input := map[string]any{"PageNumber": page, "PageSize": int32(100)}
		if cfg.project != "" {
			input["ProjectName"] = cfg.project
		}
		out := &managedEndpointsOutput{}
		req := client.NewRequest(&request.Operation{
			Name:       "InnerDescribeModelEndpoints",
			HTTPMethod: http.MethodPost,
			HTTPPath:   "/",
		}, &input, out)
		req.HTTPRequest.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.SetContext(ctx)
		if err := req.Send(); err != nil {
			return nil, skipped, err
		}
		for _, endpoint := range out.Items {
			if endpoint == nil || !strings.EqualFold(volcengine.StringValue(endpoint.Status), "Running") {
				skipped.nonRunning++
				continue
			}
			modelName := strings.TrimSpace(volcengine.StringValue(endpoint.ModelId))
			if modelName == "" {
				modelName = managedFoundationModelName(endpoint)
			}
			if modelName == "" {
				modelName = strings.TrimSpace(volcengine.StringValue(endpoint.Name))
			}
			if modelName == "" {
				skipped.uncallable++
				continue
			}
			deployments = append(deployments, desiredDeployment(
				modelName,
				volcengine.StringValue(endpoint.Id),
				modelName,
				"model",
				valueOr(volcengine.StringValue(endpoint.ProjectName), cfg.project),
				cfg.arkAPIKey,
				modelTypeFoundation,
				endpointSourceBuiltIn,
			))
		}
		if len(out.Items) == 0 || out.TotalCount != nil && page*100 >= *out.TotalCount || out.TotalCount == nil && len(out.Items) < 100 {
			break
		}
	}
	return deployments, skipped, nil
}

func managedFoundationModelName(endpoint *managedEndpoint) string {
	if endpoint == nil || endpoint.ModelReference == nil || endpoint.ModelReference.FoundationModel == nil {
		return ""
	}
	model := endpoint.ModelReference.FoundationModel
	name := strings.TrimSpace(volcengine.StringValue(model.Name))
	version := strings.TrimSpace(volcengine.StringValue(model.ModelVersion))
	if name != "" && version != "" {
		return name + "-" + version
	}
	return name
}

func (s skipStats) total() int {
	return s.nonRunning + s.uncallable + s.shadowedCustom + s.duplicate
}

func arkModelType(endpoint *ark.ItemForListEndpointsOutput) string {
	if endpoint == nil {
		return ""
	}
	switch volcengine.StringValue(endpoint.EndpointModelType) {
	case "FoundationModel":
		return modelTypeFoundation
	case "CustomModel":
		return modelTypeCustom
	case "ServiceModel":
		return modelTypeService
	}
	if endpoint.ModelReference != nil && endpoint.ModelReference.CustomModelId != nil {
		return modelTypeCustom
	}
	return modelTypeFoundation
}

func arkModelName(endpoint *ark.ItemForListEndpointsOutput, modelType string) string {
	if endpoint == nil {
		return ""
	}
	if ref := endpoint.ModelReference; ref != nil {
		if modelType == modelTypeCustom {
			if id := strings.TrimSpace(volcengine.StringValue(ref.CustomModelId)); id != "" {
				return id
			}
			if name := strings.TrimSpace(volcengine.StringValue(endpoint.Name)); name != "" {
				return name
			}
			return strings.TrimSpace(volcengine.StringValue(endpoint.Id))
		}
		if name := arkFoundationModelName(endpoint); name != "" {
			return name
		}
		if id := strings.TrimSpace(volcengine.StringValue(ref.CustomModelId)); id != "" {
			return id
		}
	}
	if name := strings.TrimSpace(volcengine.StringValue(endpoint.Name)); name != "" {
		return name
	}
	return strings.TrimSpace(volcengine.StringValue(endpoint.Id))
}

func arkFoundationModelName(endpoint *ark.ItemForListEndpointsOutput) string {
	if endpoint == nil || endpoint.ModelReference == nil || endpoint.ModelReference.FoundationModel == nil {
		return ""
	}
	model := endpoint.ModelReference.FoundationModel
	name := strings.TrimSpace(volcengine.StringValue(model.Name))
	version := strings.TrimSpace(volcengine.StringValue(model.ModelVersion))
	if name != "" && version != "" {
		return name + "-" + version
	}
	return name
}

func arkInferenceTarget(endpoint *ark.ItemForListEndpointsOutput, modelType string) (string, string) {
	if endpointID := strings.TrimSpace(volcengine.StringValue(endpoint.Id)); endpointID != "" {
		return endpointID, "endpoint"
	}
	if modelType != modelTypeCustom {
		if modelID := arkFoundationModelName(endpoint); modelID != "" {
			return modelID, "model"
		}
	}
	return "", ""
}

func desiredDeployment(modelName, endpointID, inferenceID, inferenceType, project, apiKey, modelType, endpointSource string) deployment {
	sum := sha256.Sum256([]byte(apiKey))
	return deployment{
		ModelName: modelName,
		LiteLLMParams: map[string]any{
			"model":   "volcengine/" + inferenceID,
			"api_key": apiKey,
		},
		ModelInfo: map[string]any{
			managedKey:            true,
			"ark_endpoint_id":     endpointID,
			"ark_inference_id":    inferenceID,
			"ark_inference_type":  inferenceType,
			"ark_project_name":    project,
			"ark_model_type":      modelType,
			"ark_endpoint_source": endpointSource,
			"ark_sync_key_hash":   hex.EncodeToString(sum[:]),
		},
	}
}

func normalizeDeployments(deployments []deployment) ([]deployment, int, int) {
	builtInNames := make(map[string]bool)
	for _, model := range deployments {
		if stringField(model.ModelInfo, "ark_endpoint_source") == endpointSourceBuiltIn {
			builtInNames[model.ModelName] = true
		}
	}
	filtered := deployments[:0]
	seen := make(map[string]bool)
	shadowed, duplicate := 0, 0
	for _, model := range deployments {
		if stringField(model.ModelInfo, "ark_endpoint_source") != endpointSourceBuiltIn && builtInNames[model.ModelName] {
			shadowed++
			continue
		}
		if id := managedID(model); seen[id] {
			duplicate++
			continue
		} else {
			seen[id] = true
		}
		filtered = append(filtered, model)
	}
	return filtered, shadowed, duplicate
}

func syncModels(ctx context.Context, client *liteLLMClient, desired []deployment, apply bool, out io.Writer) error {
	current, err := client.models(ctx)
	if err != nil {
		return fmt.Errorf("read LiteLLM models: %w", err)
	}
	actions := reconcile(current, desired)
	if len(actions) == 0 {
		fmt.Fprintln(out, "LiteLLM is already in sync")
		return nil
	}
	for _, action := range actions {
		fmt.Fprintf(out, "%s %s (%s)\n", strings.ToUpper(action.kind), action.model.ModelName, managedID(action.model))
	}
	if !apply {
		fmt.Fprintln(out, "preview only; rerun with --apply")
		return nil
	}

	for _, action := range actions {
		if action.kind == "delete" {
			continue
		}
		if err := client.write(ctx, action); err != nil {
			return fmt.Errorf("%s %s: %w", action.kind, action.model.ModelName, err)
		}
	}
	for _, action := range actions {
		if action.kind != "delete" {
			continue
		}
		if err := client.write(ctx, action); err != nil {
			return fmt.Errorf("delete %s: %w", action.model.ModelName, err)
		}
	}
	fmt.Fprintln(out, "sync complete")
	return nil
}

func reconcile(current, desired []deployment) []action {
	managed := make(map[string][]deployment)
	for _, model := range current {
		if boolField(model.ModelInfo, managedKey) {
			managed[managedID(model)] = append(managed[managedID(model)], model)
		}
	}

	seen := make(map[string]bool)
	var actions []action
	for _, want := range desired {
		syncID := managedID(want)
		existing := managed[syncID]
		if len(existing) == 0 {
			actions = append(actions, action{kind: "add", model: want})
			continue
		}
		have := existing[0]
		seen[stringField(have.ModelInfo, "id")] = true
		if !sameManagedDeployment(have, want) {
			actions = append(actions, action{kind: "update", model: want, id: stringField(have.ModelInfo, "id")})
		}
	}
	for _, models := range managed {
		for _, model := range models {
			id := stringField(model.ModelInfo, "id")
			if !seen[id] {
				actions = append(actions, action{kind: "delete", model: model, id: id})
			}
		}
	}
	sort.SliceStable(actions, func(i, j int) bool {
		if actions[i].kind == actions[j].kind {
			return managedID(actions[i].model) < managedID(actions[j].model)
		}
		order := map[string]int{"add": 0, "update": 1, "delete": 2}
		return order[actions[i].kind] < order[actions[j].kind]
	})
	return actions
}

func sameManagedDeployment(have, want deployment) bool {
	return have.ModelName == want.ModelName &&
		stringField(have.LiteLLMParams, "model") == stringField(want.LiteLLMParams, "model") &&
		stringField(have.ModelInfo, "ark_endpoint_id") == stringField(want.ModelInfo, "ark_endpoint_id") &&
		stringField(have.ModelInfo, "ark_inference_id") == stringField(want.ModelInfo, "ark_inference_id") &&
		stringField(have.ModelInfo, "ark_inference_type") == stringField(want.ModelInfo, "ark_inference_type") &&
		stringField(have.ModelInfo, "ark_project_name") == stringField(want.ModelInfo, "ark_project_name") &&
		stringField(have.ModelInfo, "ark_model_type") == stringField(want.ModelInfo, "ark_model_type") &&
		stringField(have.ModelInfo, "ark_endpoint_source") == stringField(want.ModelInfo, "ark_endpoint_source") &&
		stringField(have.ModelInfo, "ark_sync_key_hash") == stringField(want.ModelInfo, "ark_sync_key_hash")
}

func managedID(model deployment) string {
	return valueOr(stringField(model.ModelInfo, "ark_inference_id"), stringField(model.ModelInfo, "ark_endpoint_id"))
}

func (c *liteLLMClient) models(ctx context.Context) ([]deployment, error) {
	var response struct {
		Data []deployment `json:"data"`
	}
	if err := c.request(ctx, http.MethodGet, "/model/info", nil, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *liteLLMClient) write(ctx context.Context, action action) error {
	switch action.kind {
	case "add":
		return c.request(ctx, http.MethodPost, "/model/new", action.model, nil)
	case "update":
		model := action.model
		model.ModelInfo = cloneMap(model.ModelInfo)
		model.ModelInfo["id"] = action.id
		return c.request(ctx, http.MethodPost, "/model/update", model, nil)
	case "delete":
		if action.id == "" {
			return errors.New("managed LiteLLM model has no model_info.id")
		}
		return c.request(ctx, http.MethodPost, "/model/delete", map[string]string{"id": action.id}, nil)
	default:
		return fmt.Errorf("unknown action %q", action.kind)
	}
}

func (c *liteLLMClient) request(ctx context.Context, method, path string, body, result any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.masterKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		message := strings.ReplaceAll(string(data), c.masterKey, "[REDACTED]")
		message = strings.ReplaceAll(message, c.arkAPIKey, "[REDACTED]")
		return fmt.Errorf("LiteLLM %s %s returned %s: %s", method, path, resp.Status, strings.TrimSpace(message))
	}
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode LiteLLM response: %w", err)
		}
	}
	return nil
}

func stringField(values map[string]any, key string) string {
	value, _ := values[key].(string)
	return value
}

func boolField(values map[string]any, key string) bool {
	value, _ := values[key].(bool)
	return value
}

func cloneMap(source map[string]any) map[string]any {
	target := make(map[string]any, len(source)+1)
	for key, value := range source {
		target[key] = value
	}
	return target
}

func valueOr(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

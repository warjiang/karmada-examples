package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/volcengine/volcengine-go-sdk/service/ark"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func TestSyncModels(t *testing.T) {
	var writes []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer master" {
			t.Fatalf("missing authorization")
		}
		if r.Method == http.MethodGet && r.URL.Path == "/model/info" {
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []deployment{
				{
					ModelName:     "old-name",
					LiteLLMParams: map[string]any{"model": "volcengine/ep-1"},
					ModelInfo:     map[string]any{"id": "id-1", managedKey: true, "ark_endpoint_id": "ep-1"},
				},
				{
					ModelName:     "stale",
					LiteLLMParams: map[string]any{"model": "volcengine/ep-old"},
					ModelInfo:     map[string]any{"id": "id-old", managedKey: true, "ark_endpoint_id": "ep-old"},
				},
				{
					ModelName: "manual",
					ModelInfo: map[string]any{"id": "manual"},
				},
			}})
			return
		}
		data, _ := io.ReadAll(r.Body)
		writes = append(writes, r.URL.Path+" "+string(data))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &liteLLMClient{baseURL: server.URL, masterKey: "master", arkAPIKey: "ark-key", httpClient: server.Client()}
	desired := []deployment{
		desiredDeployment("model-v1", "ep-1", "ep-1", "endpoint", "default", "ark-key", modelTypeFoundation, endpointSourceUser),
		desiredDeployment("model-v2", "ep-2", "ep-2", "endpoint", "default", "ark-key", modelTypeCustom, endpointSourceUser),
	}
	var output strings.Builder
	if err := syncModels(context.Background(), client, desired, true, &output); err != nil {
		t.Fatal(err)
	}
	if len(writes) != 3 {
		t.Fatalf("got %d writes, want 3: %v", len(writes), writes)
	}
	if !strings.HasPrefix(writes[0], "/model/new ") ||
		!strings.HasPrefix(writes[1], "/model/update ") ||
		!strings.HasPrefix(writes[2], "/model/delete ") {
		t.Fatalf("writes must upsert before delete: %v", writes)
	}
	if strings.Contains(strings.Join(writes, "\n"), "manual") {
		t.Fatalf("manual model was modified: %v", writes)
	}
	if !strings.Contains(writes[1], `"id":"id-1"`) || !strings.Contains(writes[2], `"id":"id-old"`) {
		t.Fatalf("model IDs missing from update/delete: %v", writes)
	}
}

func TestARKModelName(t *testing.T) {
	endpoint := &ark.ItemForListEndpointsOutput{
		Id:                volcengine.String("ep-1"),
		EndpointModelType: volcengine.String("FoundationModel"),
		ModelReference: &ark.ModelReferenceForListEndpointsOutput{
			FoundationModel: &ark.FoundationModelForListEndpointsOutput{
				Name:         volcengine.String("doubao-pro"),
				ModelVersion: volcengine.String("240515"),
			},
		},
	}
	if got := arkModelName(endpoint, arkModelType(endpoint)); got != "doubao-pro-240515" {
		t.Fatalf("got %q", got)
	}
	endpoint.Id = nil
	if id, kind := arkInferenceTarget(endpoint, modelTypeFoundation); id != "doubao-pro-240515" || kind != "model" {
		t.Fatalf("got fallback %q/%q", id, kind)
	}

	endpoint.EndpointModelType = volcengine.String("CustomModel")
	endpoint.ModelReference.CustomModelId = volcengine.String("cm-user-model")
	if got := arkModelName(endpoint, arkModelType(endpoint)); got != "cm-user-model" {
		t.Fatalf("custom model got %q", got)
	}
	if id, _ := arkInferenceTarget(endpoint, modelTypeCustom); id != "" {
		t.Fatalf("custom model without endpoint ID got unsafe fallback %q", id)
	}
}

func TestBuiltInWinsNameConflict(t *testing.T) {
	models := []deployment{
		desiredDeployment("same-name", "custom-1", "custom-1", "endpoint", "default", "key", modelTypeCustom, endpointSourceUser),
		desiredDeployment("same-name", "ep-m-1", "same-name", "model", "default", "key", modelTypeFoundation, endpointSourceBuiltIn),
		desiredDeployment("same-name", "ep-m-2", "same-name", "model", "default", "key", modelTypeFoundation, endpointSourceBuiltIn),
		desiredDeployment("custom-only", "custom-2", "custom-2", "endpoint", "default", "key", modelTypeCustom, endpointSourceUser),
	}
	got, shadowed, duplicate := normalizeDeployments(models)
	if shadowed != 1 || duplicate != 1 || len(got) != 2 {
		t.Fatalf("got %d models, %d shadowed and %d duplicate", len(got), shadowed, duplicate)
	}
	for _, model := range got {
		if model.ModelName == "same-name" && stringField(model.ModelInfo, "ark_endpoint_source") != endpointSourceBuiltIn {
			t.Fatal("custom model survived a built-in name conflict")
		}
	}
}

func TestManagedEndpointUsesModelID(t *testing.T) {
	endpoint := &managedEndpoint{
		Id:      volcengine.String("ep-m-20260917111620-z6gqn"),
		ModelId: volcengine.String("deepseek-v4-1-flash-260910"),
	}
	modelID := strings.TrimSpace(volcengine.StringValue(endpoint.ModelId))
	model := desiredDeployment(
		modelID,
		volcengine.StringValue(endpoint.Id),
		modelID,
		"model",
		"default",
		"key",
		modelTypeFoundation,
		endpointSourceBuiltIn,
	)
	if managedID(model) != "deepseek-v4-1-flash-260910" ||
		stringField(model.LiteLLMParams, "model") != "volcengine/deepseek-v4-1-flash-260910" ||
		stringField(model.ModelInfo, "ark_endpoint_id") != "ep-m-20260917111620-z6gqn" ||
		stringField(model.ModelInfo, "ark_endpoint_source") != endpointSourceBuiltIn {
		t.Fatalf("unexpected managed model: %#v", model)
	}
}

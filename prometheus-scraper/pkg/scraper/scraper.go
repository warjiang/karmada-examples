package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"prometheus-scraper/pkg/database"
)

func ScrapeFromPod(namespace, name, port string) ([]*database.Metric, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/warjiang/.kube/karmada.config")
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	data, err := clientset.CoreV1().RESTClient().Get().
		Namespace(namespace).Resource("pods").
		Name(fmt.Sprintf("%s:%s", name, port)).
		SubResource("proxy").Suffix("metrics").
		Do(context.Background()).Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics from pod: %w", err)
	}

	return parseMetrics(bytes.NewReader(data))
}

func Scrape(url string) ([]*database.Metric, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to scrape metrics: status code %d", resp.StatusCode)
	}

	return parseMetrics(resp.Body)
}

func parseMetrics(reader io.Reader) ([]*database.Metric, error) {
	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metrics: %w", err)
	}

	var metrics []*database.Metric
	for _, mf := range metricFamilies {
		for _, m := range mf.GetMetric() {
			metric := &database.Metric{
				Name: mf.GetName(),
			}
			for _, lp := range m.GetLabel() {
				metric.Labels = append(metric.Labels, database.Label{
					Name:  lp.GetName(),
					Value: lp.GetValue(),
				})
			}
			switch mf.GetType() {
			case dto.MetricType_COUNTER:
				metric.Value = m.GetCounter().GetValue()
			case dto.MetricType_GAUGE:
				metric.Value = m.GetGauge().GetValue()
			case dto.MetricType_HISTOGRAM:
				metric.Value = float64(m.GetHistogram().GetSampleCount())
			case dto.MetricType_SUMMARY:
				metric.Value = float64(m.GetSummary().GetSampleCount())
			}
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

package provider

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/metrics/pkg/apis/custom_metrics"

	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider/helpers"

	"net/http"
)

type CustomMetricsProvider interface {
	ListAllMetrics() []provider.CustomMetricInfo
	GetMetricByName(name types.NamespacedName, info provider.CustomMetricInfo, labels labels.Selector) (*custom_metrics.MetricValue, error)
	GetMetricBySelector(name string, selector labels.Selector, info provider.CustomMetricInfo) (*custom_metrics.MetricValueList, error)
}

type Instances struct {
	Count int `json:"count"`
}

func (p *camundaProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return []provider.CustomMetricInfo{
		{
			GroupResource: schema.GroupResource{Group: "", Resource: "services"},
			Metric:        "camunda_queue_count",
			Namespaced:    true,
		},
	}
}

type camundaProvider struct {
	client dynamic.Interface
	mapper apimeta.RESTMapper

	// just increment values when they're requested
	values map[provider.CustomMetricInfo]int64
}

func getCamundaProcesses() int {
	camundaApiUrl := "http://camunda-service.default.svc.cluster.local:8080/engine-rest/history/process-instance/count"
	formatedDate := time.Now().Add(-10 * time.Second).Format("2006-01-02T15:04:05.000-0700")
	jsonStr := []byte(`{"startedAfter":"` + formatedDate + `"}`)

	req, err := http.NewRequest("POST", camundaApiUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var instances Instances
	err = json.Unmarshal(body, &instances)
	return instances.Count
}

func (p *camundaProvider) GetMetricByName(name types.NamespacedName, info provider.CustomMetricInfo, label labels.Selector) (*custom_metrics.MetricValue, error) {
	objRef, err := helpers.ReferenceFor(p.mapper, name, info)
	if err != nil {
		return nil, err
	}

	// Camunda-HPA
	processes_started := getCamundaProcesses()

	return &custom_metrics.MetricValue{
		DescribedObject: objRef,
		// you'll want to use the actual timestamp in a real adapter
		Timestamp: metav1.Time{time.Now()},
		Value:     *resource.NewQuantity(int64(processes_started), resource.DecimalSI),
	}, nil
}

func (p *camundaProvider) GetMetricBySelector(namespace string, selector labels.Selector, info provider.CustomMetricInfo, labels labels.Selector) (*custom_metrics.MetricValueList, error) {
	return nil, nil
}

func NewProvider(client dynamic.Interface, mapper apimeta.RESTMapper) provider.CustomMetricsProvider {
	provider := &camundaProvider{
		client: client,
		mapper: mapper,
		values: make(map[provider.CustomMetricInfo]int64),
	}
	return provider
}

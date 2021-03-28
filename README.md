# Camunda k8s Adapter
Provides metrics to Kubernetes via custom *custom.metrics.k8s.io* api.

The adapter has a provider which query Camunda API for the number of processes started on the last 10s and return it to Kubernetes.

Based on this metric, Kubernetes HPA (Horizontal Pod Autoscaler) can be configured.

# Deploy
Just ran the Kubernetes resources manifest:
`kubectl apply -f camunda-k8s-adapter.yaml`

# Notes
- The resources will be deployed in custom-metrics namespace.
- The Camunda API should be available on this URL:
http://camunda-service.default.svc.cluster.local:8080/engine-rest/ 

# How to test
Check *custom.metrics.k8s.io/v1beta1* API to verify the metrics are available:
```
➜  ~ kubectl get --raw="/apis/custom.metrics.k8s.io/v1beta1/namespaces/kube-system/services/custom-metrics/camunda_queue_count" | jq
{
  "kind": "MetricValueList",
  "apiVersion": "custom.metrics.k8s.io/v1beta1",
  "metadata": {
    "selfLink": "/apis/custom.metrics.k8s.io/v1beta1/namespaces/kube-system/services/custom-metrics/camunda_queue_count"
  },
  "items": [
    {
      "describedObject": {
        "kind": "Service",
        "namespace": "kube-system",
        "name": "custom-metrics",
        "apiVersion": "/v1"
      },
      "metricName": "",
      "timestamp": "2021-03-28T16:15:00Z",
      "value": "50",
      "selector": null
    }
  ]
}
```
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
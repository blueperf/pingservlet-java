kubectl delete -f deploy/crds/benchmark_v1alpha1_pingservlet_cr.yaml
sleep 5
kubectl delete -f deploy/operator.yaml

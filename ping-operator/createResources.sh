kubectl apply -f deploy/operator.yaml
sleep 20
kubectl apply -f deploy/crds/benchmark_v1alpha1_pingservlet_cr.yaml

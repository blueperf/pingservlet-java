# PingServlet Application
This application is a simple Java servlet application that prints out the total number of hits of the page.  It utilizes binary output for faster response time (it is about 8x faster than PrintWriter) for simplicity.  This is implemented by using [Operator SDK](https://github.com/operator-framework/operator-sdk) to natively create the PingServlet Custom Resources into Kubernetes environment.  It is used to focus on ease of maintenance by using "kubectl" command.  Also, it utilizes Kubernetes default rolling update to achieve zero downtime. 

## Configuration

Please use ping-operator/deploy/crds/benchmark_v1alpha1_pingservlet_cr.yaml as your template for your Custom Resources (CR).  User is able to modify any values except the following:

  - apiVersion: benchmark.perf/v1alpha1
  - kind: PingServlet
  - host: (hostname is final after the initial deployment)
  
## Installation

PingServlet Operator will be installed by running following commands (all these resources are needed to create the operator):

 - Setup Service Account
 
```console
$ kubectl apply -f ./ping-operator/deploy/service_account.yaml
```
 - Setup RBAC
 
```console
$ kubectl apply -f ./ping-operator/deploy/role.yaml
$ kubectl apply -f ./ping-operator/deploy/role_binding.yaml
```
 - Setup Custom Resource Definition (CRD)
 
```console
$ kubectl apply -f ./ping-operator/deploy/crds/benchmark_v1alpha1_pingservlet_crd.yaml
```
 - Deploy ping-operator
 
```console
$ kubectl apply -f ./ping-operator/deploy/operator.yaml
```

PingServlet Instances will be deployed by the PingServlet Operator by running the following command:

```console
$ kubectl apply -f ./ping-operator/deploy/crds/benchmark_v1alpha1_pingservlet_cr.yaml
```

## Update PingServlet instances

To update the deployed PingServlet services version, update the image name in the cr.yaml file (e.g. ping-operator/deploy/crds/benchmark_v1alpha1_pingservlet_cr.yaml), then run the following command (use your cr.yaml file path). It will use Kubernetes default Rolling Update to update each instances without any down time.

```console
$ kubectl apply -f ./ping-operator/deploy/crds/benchmark_v1alpha1_pingservlet_cr.yaml
```
Any update in the cr.yaml file will be applied by above command (e.g. changing maxReplicas, targetCPUPercent, totalUsers)


## URL
This is the URL to access the application

```console
http://<host>:<port>/servlet/PingServlet
```

(e.g. http://ping.your_cluster_name.us-south.containers.appdomain.cloud/bff)

## jMeter load test
Please use the included PingServlet-v1.jmx file to run the test

```console
jmeter  -n -t PingServlet-v1.jmx -j ping.log -JTHREAD=1 -JDURATION=60 -JRAMP=0 -JPORT=${PORT} -JHOST=${HOST} ;
```

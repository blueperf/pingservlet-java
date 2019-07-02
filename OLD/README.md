# pingservlet-java - NOT SUPPORTED (Provided AS IS)
This application is a simple Java servlet application that prints out the total number of hits of the page.  It utilizes binary output for faster response time (it is about 8x faster than PrintWriter) for simplicity. 


## Configuration

The following table lists the configurable parameters of this chart and their default values.
The parameters allow you to:
* change the image of any microservice from the one provided by IBM to one that you build (e.g. if you want to try to modify a service)
* change the resource utilization and autoscale CPU percentage.

| Parameter                           | Description                                         | Default                                                                         |
| ----------------------------------- | ----------------------------------------------------| --------------------------------------------------------------------------------|
| | | |
| image.repository | image repository |  watsoncloudperf/pingservlet-java |
| image.tag | image tag |  latest |
| resources.requests.cpu | CPU resource request | 100m |
| resources.requests.memory | memory resource request | 128Mi |
| hpa.targetCPUUtilizationPercentage | CPU utilization percentage for autoscale | 80 |
| ingress.host | ingress host name |  |


## Installing the Chart

You can install the chart by running the following command (add correct ingress subdomain. above parameters can be dynamically set by --set option):

```console
helm install pingservlet-java --name pingservlet --set ingress.host=pingservlet.<subdomain>
```

## URL
This is the URL to access the application 

```console
http://<host>:<port>/servlet/PingServlet 
```

(e.g. http://ping.us-south.containers.mybluemix.net/servlet/PingServlet)

## Optional

Other additional scripts are provided for the convenience for creating & deploying your own image into different kubernetes environments.

  - src files
  - pom.xml (building war file using maven)
  - war file
  - server.xml (for Liberty profile)
  - Dockerfile_KS (Creating docker image) 
  - various yaml files (Deploying the docker image to various Kubernetes environments)


## jMeter load test
Please use the included PingServlet-v1.jmx file to run the test

```console
jmeter  -n -t PingServlet-v1.jmx -j ping.log -JTHREAD=1 -JDURATION=60 -JRAMP=0 -JPORT=${PORT} -JHOST=${HOST} ;
```

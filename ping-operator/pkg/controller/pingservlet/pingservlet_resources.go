package pingservlet

import (
	"fmt"

	benchmarkv1alpha1 "ping-operator/pkg/apis/benchmark/v1alpha1"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func labelsForPingServlet(name string) map[string]string {
	return map[string]string{"app": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func deploymentForPingServlet(cr *benchmarkv1alpha1.PingServlet) *extensionsv1beta1.Deployment {
	replicas := cr.Spec.Size
	dep := &extensionsv1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labelsForPingServlet(cr.Name),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelsForPingServlet(cr.Name),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: cr.Spec.Image,
						Name:  cr.Name,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								"cpu": resource.MustParse("700m"),
							},
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/servlet/PingServlet",
									Port: intstr.FromInt(int(cr.Spec.Port)),
								},
							},
							InitialDelaySeconds: 10,
							TimeoutSeconds:      5,
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/servlet/PingServlet",
									Port: intstr.FromInt(int(cr.Spec.Port)),
								},
							},
							InitialDelaySeconds: 120,
						},
					}},
				},
			},
		},
	}
	return dep
}

func serviceForPingServlet(cr *benchmarkv1alpha1.PingServlet) *corev1.Service {
	annotations := map[string]string{
		"description": "PingServlet service",
	}
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core/v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-svc", cr.Name),
			Namespace:   cr.Namespace,
			Labels:      labelsForPingServlet(cr.Name),
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceType("ClusterIP"),
			Selector: labelsForPingServlet(cr.Name),
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     int32(cr.Spec.Port),
					Protocol: "TCP",
				},
			},
		},
	}
	return service
}

func ingressForPingServlet(cr *benchmarkv1alpha1.PingServlet) *extensionsv1beta1.Ingress {
	annotations := map[string]string{
		"description": "PingServlet Ingress",
	}
	ing := &extensionsv1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-ing", cr.Name),
			Namespace:   cr.Namespace,
			Labels:      labelsForPingServlet(cr.Name),
			Annotations: annotations,
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{{
				Host: cr.Spec.Host,
				IngressRuleValue: extensionsv1beta1.IngressRuleValue{
					HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
						Paths: []extensionsv1beta1.HTTPIngressPath{{
							Path: "/servlet",
							Backend: extensionsv1beta1.IngressBackend{
								ServiceName: fmt.Sprintf("%s-svc", cr.Name),
								ServicePort: intstr.FromInt(int(cr.Spec.Port)),
							},
						}},
					},
				},
			}},
		},
	}
	return ing
}

func hpaForPingServlet(cr *benchmarkv1alpha1.PingServlet) *autoscalingv1.HorizontalPodAutoscaler {
	annotations := map[string]string{
		"description": "PingServlet service",
	}
	hpa := &autoscalingv1.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "autoscaling/v1",
			Kind:       "HorizontalPodAutoscaler",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-hpa", cr.Name),
			Namespace:   cr.Namespace,
			Labels:      labelsForPingServlet(cr.Name),
			Annotations: annotations,
		},
		Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
			MinReplicas:                    &cr.Spec.MinReplicas,
			MaxReplicas:                    cr.Spec.MaxReplicas,
			TargetCPUUtilizationPercentage: &cr.Spec.TargetCPUPercent,
			ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       cr.Name,
				APIVersion: "extensions/v1beta1",
			},
		},
	}
	return hpa
}

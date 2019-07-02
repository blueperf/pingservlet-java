package pingservlet

import (
	"context"
	"fmt"
	"reflect"

	benchmarkv1alpha1 "ping-operator/pkg/apis/benchmark/v1alpha1"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_pingservlet")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new PingServlet Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePingServlet{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("pingservlet-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource PingServlet
	err = c.Watch(&source.Kind{Type: &benchmarkv1alpha1.PingServlet{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Deployment and requeue the owner Olb
	err = c.Watch(&source.Kind{Type: &extensionsv1beta1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &benchmarkv1alpha1.PingServlet{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcilePingServlet{}

// ReconcilePingServlet reconciles a PingServlet object
type ReconcilePingServlet struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a PingServlet object and makes changes based on the state read
// and what is in the PingServlet.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePingServlet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PingServlet")

	// Fetch the PingServlet instance
	instance := &benchmarkv1alpha1.PingServlet{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	/*************************************
	Create Ingress Resource
	**************************************/

	//Define a new resource object
	ing := ingressForPingServlet(instance)

	//Set instance as the owner of the service
	if err := controllerutil.SetControllerReference(instance, ing, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Check if this Ingress already exists
	ingress := &extensionsv1beta1.Ingress{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ing.Name, Namespace: ing.Namespace}, ingress)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.V(2).Info("Creating a new Ingress", "ing.Namespace", ing.Namespace, "ing.Name", ing.Name)
		err = r.client.Create(context.TODO(), ing)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Ingress", "ing.Namespace", ing.Namespace, "ing.Name", ing.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Ingress")
		return reconcile.Result{}, err
	}

	/*************************************
	Create Serevice Resource
	**************************************/
	reqLogger.V(1).Info("Creating a Services")

	svc := serviceForPingServlet(instance)

	//Set instance as the owner of the service
	if err := controllerutil.SetControllerReference(instance, svc, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Service already exists
	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.V(2).Info("Creating a new Service", "svc.Namespace", svc.Namespace, "svc.Name", svc.Name)
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "svc.Namespace", svc.Namespace, "svc.Name", svc.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}

	/*************************************
	Create HPA Resource
	**************************************/
	reqLogger.V(1).Info("Creating a HPA Resources")

	//Define a new resource object
	hpa := hpaForPingServlet(instance)

	//Set instance as the owner of the service
	if err := controllerutil.SetControllerReference(instance, hpa, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Check if this HPA already exists
	horizontalPodAutoscaler := &autoscalingv1.HorizontalPodAutoscaler{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: hpa.Name, Namespace: hpa.Namespace}, horizontalPodAutoscaler)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.V(2).Info("Creating a new HorizontalPodAutoscaler", "hpa.Namespace", hpa.Namespace, "hpa.Name", hpa.Name)
		err = r.client.Create(context.TODO(), hpa)
		if err != nil {
			reqLogger.Error(err, "Failed to create new HorizontalPodAutoscaler", "hpa.Namespace", hpa.Namespace, "hpa.Name", hpa.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get HorizontalPodAutoscaler")
		return reconcile.Result{}, err
	}

	/*************************************
	Create Deplyment Resource
	**************************************/

	reqLogger.V(1).Info("Creating a Deployments")

	//Define a new resource object
	dep := deploymentForPingServlet(instance)

	//Set instance as the owner of the service
	if err := controllerutil.SetControllerReference(instance, dep, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Check if this Deployment already exists
	deployment := &extensionsv1beta1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		reqLogger.V(2).Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}

		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	/*************************************
	Update Resources
	**************************************/

	reqLogger.V(1).Info("Update Resources")

	/******** Update Deployment Size *******
	UNCOMMENT TO MANUALLY ADJUST DEPLOYMENT SIZE INSTEAD OF USING HPA
	size := instance.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		// Update deployment
		err = r.client.Update(context.TODO(), deployment)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment Size", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}
	*/

	/******* Update Deployment Image *******/
	image := instance.Spec.Image
	port := int(instance.Spec.Port)
	intorstringPort := intstr.FromInt(port)
	dep = deploymentForPingServlet(instance)
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, deployment)
	reqLogger.V(2).Info(fmt.Sprintf("Getting deployment for deployment %s", deployment.Name))
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Error(err, fmt.Sprintf("Deployment %s is not found", dep.Name))
	} else if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to get Deployment %s", dep.Name))
	}
	if err == nil {
		if deployment.Spec.Template.Spec.Containers[0].Image != image || deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Port != intorstringPort || deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Port != intorstringPort {
			if deployment.Spec.Template.Spec.Containers[0].Image != image {
				deployment.Spec.Template.Spec.Containers[0].Image = image
			}
			if deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Port != intorstringPort {
				deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Port = intorstringPort
			}
			if deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Port != intorstringPort {
				deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Port = intorstringPort
			}
			// Update deployment
			err = r.client.Update(context.TODO(), deployment)
			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment Image", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			reqLogger.V(2).Info(fmt.Sprintf("%s deployment is updated!", deployment.Name))
			return reconcile.Result{Requeue: true}, nil
		}
	}

	/******* Update Service *******/
	svc = serviceForPingServlet(instance)
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, service)
	reqLogger.V(2).Info(fmt.Sprintf("Getting service for svc %s", service.Name))
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Error(err, fmt.Sprintf("Service %s is not found", svc.Name))
	} else if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to get Service %s", svc.Name))
	}
	if err == nil {
		if service.Spec.Ports[0].Port != int32(port) {
			service.Spec.Ports[0].Port = int32(port)
			// Update service
			err = r.client.Update(context.TODO(), service)
			if err != nil {
				reqLogger.Error(err, "Failed to update Service port", "Service.Namespace", service.Namespace, "service.Name", service.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			reqLogger.V(2).Info(fmt.Sprintf("%s service is updated!", service.Name))
			return reconcile.Result{Requeue: true}, nil
		}
	}

	reqLogger.V(2).Info(fmt.Sprintf("%s deployment HPA check", deployment.Name))

	/******* Update HPA variables *******/
	min := &instance.Spec.MinReplicas
	max := instance.Spec.MaxReplicas
	cpu := &instance.Spec.TargetCPUPercent

	hpa = hpaForPingServlet(instance)
	// Check if this HPA already exists
	horizontalPodAutoscaler = &autoscalingv1.HorizontalPodAutoscaler{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: hpa.Name, Namespace: hpa.Namespace}, horizontalPodAutoscaler)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Error(err, fmt.Sprintf("HPA %s is not found", hpa.Name))
	} else if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to get HPA %s", hpa.Name))
	} else {
		if *horizontalPodAutoscaler.Spec.MinReplicas != *min || horizontalPodAutoscaler.Spec.MaxReplicas != max || *horizontalPodAutoscaler.Spec.TargetCPUUtilizationPercentage != *cpu {
			if *horizontalPodAutoscaler.Spec.MinReplicas != *min {
				*horizontalPodAutoscaler.Spec.MinReplicas = *min
			}
			if horizontalPodAutoscaler.Spec.MaxReplicas != max {
				horizontalPodAutoscaler.Spec.MaxReplicas = max
			}
			if *horizontalPodAutoscaler.Spec.TargetCPUUtilizationPercentage != *cpu {
				*horizontalPodAutoscaler.Spec.TargetCPUUtilizationPercentage = *cpu
			}

			// Update HPA
			err = r.client.Update(context.TODO(), horizontalPodAutoscaler)
			if err != nil {
				reqLogger.Error(err, "Failed to update horizontalPodAutoscaler min replicas", "horizontalPodAutoscaler.Namespace", horizontalPodAutoscaler.Namespace, "horizontalPodAutoscaler.Name", horizontalPodAutoscaler.Name)
				return reconcile.Result{}, err
			}
			reqLogger.V(2).Info(fmt.Sprintf("%s HPA is updated!", horizontalPodAutoscaler.Name))
			return reconcile.Result{Requeue: true}, nil
		}
	}

	/*************************************
	Update the deployment status with the pod names
	**************************************/
	// List the pods for this deployment's deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForPingServlet(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Application status")
			return reconcile.Result{}, err
		}
	}

	/*************************************
	END
	**************************************/

	return reconcile.Result{}, nil
}

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	//"strconv"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	closedlooppoocv1 "closedloop/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	mytools "closedloop/tools"
	networkingv1 "k8s.io/api/networking/v1"
)

// Monitoringv2Reconciler reconciles a Monitoringv2 object
type Monitoringv2Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloops,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloops/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=closedloops/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitorings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitorings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitorings/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisions/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executions/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=monitoringv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=executionv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=closedlooppooc.closedloop.io,resources=decisionv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="apps",resources=deployments/status,verbs=get;watch;list
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services/status,verbs=get;watch;list
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses/status,verbs=get;watch;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Monitoringv2 object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *Monitoringv2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	verbosityLog := new(mytools.VerbosityLog)
	verbosityLog.SetMaxLevel(2)
	l := verbosityLog.FromContext(ctx)
	l.V(1).Info("                                                                                          ")
	l.V(1).Info(">>>>>>>>>>>>>>>>>>>>>>>>>")
	l.V(1).Info("Enter Reconcile Monitoring")

	//Retreiving ClosedLoop Object who triggered the Reconciler
	nameCR := "Monitoring"
	kindCR := "Monitoringv2"
	clcm := mytools.Cr_type(nameCR, kindCR)

	//Monitoring := &closedlooppoocv1.Monitoringv2{}
	client := r.Client
	err := client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, clcm)

	v := mytools.MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}
	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it
	fmt.Printf("%s.Spec.Policy.Data %v\n", nameCR,
		v.FieldByName("Spec").FieldByName("Policy").FieldByName("Data"))

	if v.FieldByName("Name") != v.FieldByName("Status").FieldByName("Affix") {
		v.FieldByName("Status").FieldByName("Affix").Set(v.FieldByName("Name"))

		if err := r.Status().Update(ctx, clcm); err != nil {
			l.Error(err, "Failed to update "+nameCR+" status")
			return ctrl.Result{}, err
		}
		l.V(2).Info("Enter Reconcile View Spec & Status", "spec", v.FieldByName("Spec"), "status", v.FieldByName("Status"))
	}

	// If the Monitoringv2 who triggered is not find it means that it's been deleted
	if err != nil {

		// The Monitoringv2 object has been deleted or is not found, so we should delete the associated Ressource (if they exist)
		l.V(2).Info("Monitoring Instance not found, Deletion Close Loop ressources")
		if err := r.deleteDEPLOYMENT(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "Failed to delete DEPLOYMENT for Monitoring Type 2")
			return ctrl.Result{}, err
		}

		if err := r.deleteSERVICE(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "Failed to delete SERVICE for Monitoring Type 2")
			return ctrl.Result{}, err
		}

		if err := r.deleteINGRESS(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "Failed to delete INGRESS for Monitoring Type 2")
			return ctrl.Result{}, err
		}
	}

	// Creation of Ressource implementation layer (Proxy Pod, Services and Ingress) ressources  are created
	// These function are called Each time if RequestedPod=true , if there are already created nothing will happen
	if v.XFieldByName("Spec").XFieldByName("RequestedPod").Kind() != reflect.Invalid {

		if err := r.createDEPLOYMENT(ctx, clcm, l); err != nil {
			l.Error(err, "No deployment found, no deletions have been made")

		}

		if err := r.createSERVICE(ctx, clcm, l); err != nil {
			l.Error(err, "No service found, no deletions have been made")

		}

		if err := r.createINGRESS(ctx, clcm, l); err != nil {
			l.Error(err, "No ingress found, no deletions have been made")

		}

	}

	/* ---------------------------------- START Monitoring Part ---------------------------------- */

	if v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Kind").String() == "opa" {
		println("*********************************************************opa**************************")
	}

	if !mytools.OpaPolicyFoundPolicy(v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Url").String(),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Api").XFieldByName("Policy").String(),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Name").String(), l) {
		mytools.OpaPolicyWritePolicy(v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Url").String(),
			v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Api").XFieldByName("Policy").String(),
			v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Name").String(),
			v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Rule").XFieldByName("Body").String(), l)
	}
	if !mytools.OpaPolicyFoundData(v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Url").String(),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Api").XFieldByName("Data").String(),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Data").XFieldByName("Name").String(), l) {
		mytools.OpaPolicyWriteData(v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Url").String(),
			v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Api").XFieldByName("Data").String(),
			v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Data").XFieldByName("Name").String(),
			v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Data").XFieldByName("Body").String(), l)
	}

	messageinput := strings.Replace(v.XFieldByName("Spec").XFieldByName("Message").String(), "'", "\"", -1)
	l.V(1).Info("~~~~~~~~~~messageinput:  " + messageinput)
	messagecolected := mytools.MessageDataCollect(client, ctx, req,
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Input").XFieldByName("Value").String(),
		messageinput, nameCR, kindCR, l)
	l.V(1).Info("~~~~~~~~~~messagecolected:  " + messagecolected)
	message := mytools.MessageInputConvert(v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Input").XFieldByName("Schema").String(),
		messagecolected,
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Input").XFieldByName("Value").String(), l)
	l.V(1).Info("~~~~~~~~~~message:  " + message)
	answer := mytools.OpaPolicy(v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Url").String(),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Api").XFieldByName("Data").String(),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Name").String(),
		message, l)
	l.V(1).Info("~~~~~~~~~~answer:  " + answer)
	jsonData := []byte(answer)
	key := "result"
	answerinput, err := mytools.ExtractSubJSON(jsonData, key)
	if err != nil {
		fmt.Println("Error:", err)
		return ctrl.Result{}, err
	}
	l.V(1).Info("~~~~~~~~~~answerinput:  " + string(answerinput))
	messages := mytools.MessageOutputConvert(client, ctx, req,
		nameCR, kindCR,
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Result").XFieldByName("Schema").String(),
		string(answerinput),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Result").XFieldByName("Value").String(), l)

	fmt.Printf("~~~~~~~~~~messages from MessageOutputConvert:  %v", messages)

	for _, mess := range messages {
		if err := mytools.ApplyToCR(client, ctx, req, mess, l); err != nil {
			l.Error(err, "Failed to ApplyToCR")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *Monitoringv2Reconciler) deleteDEPLOYMENT(ctx context.Context, name, namespace string, l mytools.VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, deployment)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, deployment)
	if err != nil {
		l.Error(err, "Failed to delete DEPLOYMENT")
		return err
	}

	l.Info("Deleted DEPLOYMENT")
	return nil
}

func (r *Monitoringv2Reconciler) deleteSERVICE(ctx context.Context, name, namespace string, l mytools.VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	service := &v1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: name + "-deployment-service", Namespace: namespace}, service)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, service)
	if err != nil {
		l.Error(err, "Failed to delete SERVICE")
		return err
	}

	l.Info("Deleted DEPLOYMENT")
	return nil
}

func (r *Monitoringv2Reconciler) deleteINGRESS(ctx context.Context, name, namespace string, l mytools.VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	ingress := &networkingv1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: name + "-ingress", Namespace: namespace}, ingress)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, ingress)
	if err != nil {
		l.Error(err, "Failed to delete INGRESS")
		return err
	}

	l.Info("Deleted INGRESS")
	//lock.Unlock()
	return nil
}

// Function use to create the Deployment
func (r *Monitoringv2Reconciler) createDEPLOYMENT(ctx context.Context, clcm client.Object, l mytools.VerbosityLog) error {
	// Try to retrieve the Ressource to see if it's already within the Cluster
	deployment := &appsv1.Deployment{}
	v := mytools.MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}
	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it

	err := r.Get(ctx, types.NamespacedName{Name: v.FieldByName("Name").String(),
		Namespace: v.FieldByName("Namespace").String()}, deployment)

	if err == nil {
		l.V(2).Info("DEPLOYMENT Found - No Creation")
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	l.V(2).Info("DEPLOYMENT Not found")

	var replicas *int32
	var replicas_value int32 = 1

	replicas = &replicas_value

	jsonData, err := json.Marshal(v.FieldByName("Spec").FieldByName("Policy").String())
	jsonString := string(jsonData)
	fmt.Println(jsonString)

	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		jsonString = "No Data - Error"
	}
	//Creating the Deployment Object with the right Spec
	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: v.FieldByName("Namespace").String(),
			Name:      v.FieldByName("Name").String(),
		},
		Spec: appsv1.DeploymentSpec{

			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": v.FieldByName("Name").String() + "-pod",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      v.FieldByName("Namespace").String(),
					Namespace: v.FieldByName("Namespace").String(),
					Labels: map[string]string{
						"app": v.FieldByName("Name").String() + "-pod",
					},
				},

				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "restpod-receiver",
							Image: "restpod:latest",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "CLOSED_LOOP_MONITOR_NAME",
									Value: v.FieldByName("Name").String(),
								},
								{
									Name:  "CLOSED_LOOP_MONITOR_NAMESPACE",
									Value: v.FieldByName("Namespace").String(),
								},
								{
									Name:  "CLOSED_LOOP_DATA_TO_MONITOR",
									Value: jsonString,
								},
								{
									Name:  "CLOSED_LOOP_MONITOR_KIND",
									Value: v.FieldByName("Kind").String(),
								},
								{
									Name:  "CLOSED_LOOP_SOURCE_CONFIGMAP",
									Value: "sourceconfigmap",
								},
								{
									Name:  "CLOSED_LOOP_MONITORED_INTERVAL",
									Value: "30",
								},
							},
							ImagePullPolicy: v1.PullNever,
						},
					},
				},
			},
		},
	}

	l.V(2).Info("Creating DEPLOYMENT")
	return r.Create(ctx, deployment)

}

// Function use to create the Service
func (r *Monitoringv2Reconciler) createSERVICE(ctx context.Context, clcm client.Object, l mytools.VerbosityLog) error {
	// Try to retrieve the Ressource to see if it's already within the Cluster
	service := &v1.Service{}
	v := mytools.MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}
	err := r.Get(ctx, types.NamespacedName{Name: v.FieldByName("Name").String() + "-deployment-service", Namespace: v.FieldByName("Namespace").String()}, service)

	if err == nil {
		l.V(2).Info("SERVICE Found - No Creation")
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	l.V(2).Info("SERVICE Not found")
	//Creating the Service Object with the right Spec
	service = &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.FieldByName("Name").String() + "-deployment-service",
			Namespace: v.FieldByName("Namespace").String(),
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   "TCP",
				},
			},
			Selector: map[string]string{
				"app": v.FieldByName("Name").String() + "-pod",
			},
		},
	}

	l.Info("Creating DEPLOYMENT")
	return r.Create(ctx, service)

}

// Function use to create the Ingress
func (r *Monitoringv2Reconciler) createINGRESS(ctx context.Context, clcm client.Object, l mytools.VerbosityLog) error {
	// Try to retrieve the Ressource to see if it's already within the Cluster
	v := mytools.MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}
	ingress := &networkingv1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: v.FieldByName("Name").String() + "-ingress", Namespace: v.FieldByName("Namespace").String()}, ingress)

	if err == nil {
		l.V(2).Info("Ingress Found - No Creation")
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	l.Info("Ingress Not found")
	//Creating the Ingress Object with the right Spec
	ingress = &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.FieldByName("Name").String() + "-ingress",
			Namespace: v.FieldByName("Namespace").String(),
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: v.FieldByName("Name").String() + "-deployment-service.com", // Remplacez par votre nom de domaine ou adresse IP
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: (*networkingv1.PathType)(pointer.StringPtr(string(networkingv1.PathTypePrefix))),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: v.FieldByName("Name").String() + "-deployment-service",
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	l.Info("Creating Ingress")
	return r.Create(ctx, ingress)

}

func (r *Monitoringv2Reconciler) ApplyDecision(ctx context.Context, clcm client.Object, l mytools.VerbosityLog, Message string, stop bool) error {

	//	if monitoring.Spec.Time == monitoring.Spec.MonitoringPolicies.Time {
	//		return nil
	//	}
	if stop {
		return nil
	}
	v := mytools.MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}
	Decision := &closedlooppoocv1.Decision{}
	r.Get(ctx, types.NamespacedName{Name: v.FieldByName("Spec").FieldByName("Affix").String() + "-decision", Namespace: v.FieldByName("Namespace").String()}, Decision)

	l.V(1).Info("Send Message to Decision: " + "Message=" + fmt.Sprintf("%v", Message) + " Time=" + v.FieldByName("Spec").FieldByName("Time").String())

	Decision.Spec.Message = Message
	Decision.Spec.Time = v.FieldByName("Spec").FieldByName("Time").String()

	//return nil
	return r.Update(ctx, Decision)

}

// SetupWithManager sets up the controller with the Manager.
func (r *Monitoringv2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.Monitoringv2{}).
		Complete(r)
}

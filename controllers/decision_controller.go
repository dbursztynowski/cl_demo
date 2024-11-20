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
	//"encoding/json"
	"fmt"
	"reflect"

	//"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	closedlooppoocv1 "closedloop/api/v1"
	mytools "closedloop/tools"
)

// DecisionReconciler reconciles a Decision object
type DecisionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Here you give all the permission your controller need to be able to work

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
//+kubebuilder:rbac:groups="",resources=services/status,verbs=get;watch;list
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses/status,verbs=get;watch;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Decision object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DecisionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	verbosityLog := new(mytools.VerbosityLog)
	verbosityLog.SetMaxLevel(2)
	l := verbosityLog.FromContext(ctx)
	l.V(2).Info("\n********************************")
	l.V(1).Info("Enter Reconcile Decision")

	//Retreiving Decision Object who triggered the Reconciler
	nameCR := "Decision"
	kindCR := ""
	clcm := mytools.Cr_type(nameCR, kindCR)
	client := r.Client
	err := client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, clcm)
	if err != nil {
		return ctrl.Result{}, err
	}

	//r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, Decision)

	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it

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

	// START of data processing in the "message" field

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
	messagecolected := mytools.MessageDataCollect(client, ctx, req,
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Input").XFieldByName("Value").String(),
		messageinput, nameCR, kindCR, l)

	message := mytools.MessageInputConvert(v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Input").XFieldByName("Schema").String(),
		messagecolected,
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Input").XFieldByName("Value").String(), l)

	answer := mytools.OpaPolicy(v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Url").String(),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Engine").XFieldByName("Api").XFieldByName("Data").String(),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Name").String(),
		message, l)

	jsonData := []byte(answer)
	key := "result"
	answerinput, err := mytools.ExtractSubJSON(jsonData, key)
	if err != nil {
		fmt.Println("Error:", err)
		return ctrl.Result{}, err
	}
	messages := mytools.MessageOutputConvert(client, ctx, req,
		nameCR, kindCR,
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Result").XFieldByName("Schema").String(),
		string(answerinput),
		v.XFieldByName("Spec").XFieldByName("Policy").XFieldByName("Result").XFieldByName("Value").String(), l)

	for _, mess := range messages {
		if err := mytools.ApplyToCR(client, ctx, req, mess, l); err != nil {
			l.Error(err, "Failed to ApplyToCR")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// Function to update Execution Action
func (r *DecisionReconciler) ApplyExecution(ctx context.Context, clcm client.Object, l mytools.VerbosityLog, action string) error {
	// Try to retrieve the CR that we want to update

	//	if decision.Spec.Time == decision.Spec.DecisionPolicies.PrioritySpec.Time {
	//		return nil
	//	}
	v := mytools.MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}
	Execution := &closedlooppoocv1.Execution{}
	r.Get(ctx, types.NamespacedName{Name: v.FieldByName("Spec").FieldByName("Affix").String() + "-execution", Namespace: v.FieldByName("Namespace").String()}, Execution)

	l.V(2).Info("Update Action on Execution")
	//Update it's field with the variable message
	Execution.Spec.Message = action
	Execution.Spec.Time = v.FieldByName("Spec").FieldByName("Time").String()
	l.V(1).Info("Send message to Execution: " + fmt.Sprintf("%v", Execution.Spec.Message) + " " + Execution.Spec.Time)

	return r.Update(ctx, Execution)

}

func (r *DecisionReconciler) InformDeliberate(ctx context.Context, decision *closedlooppoocv1.Decision, l mytools.VerbosityLog, Message string) error {
	MonitoringDv2 := &closedlooppoocv1.MonitoringDv2{}
	r.Get(ctx, types.NamespacedName{Name: "closedloopd-v2-monitoringd", Namespace: decision.Namespace}, MonitoringDv2)

	if MonitoringDv2.Spec.Time != decision.Spec.Time {
		l.V(1).Info("Update Message on MonitoringDv2 " + MonitoringDv2.Name)
		l.V(2).Info("Data before " + fmt.Sprint(MonitoringDv2.Spec.Data))
		m := make(map[string]string)
		m["metric"] = Message
		MonitoringDv2.Spec.Data = m
		l.V(1).Info("Message send to ClosedLoop " + fmt.Sprint(MonitoringDv2.Spec.Data))
		MonitoringDv2.Spec.Time = decision.Spec.Time
		return r.Update(ctx, MonitoringDv2)
	}
	return nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *DecisionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.Decision{}).
		Complete(r)
}

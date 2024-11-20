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
Version 2
*/

package controllers

import (
	"context"
	"reflect"
	"strconv"
	"strings"

	//"sort"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	closedlooppoocv1 "closedloop/api/v1"
	mytools "closedloop/tools"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClosedLoopReconciler reconciles a ClosedLoop object
type ClosedLoopReconciler struct {
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
// the ClosedLoop object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ClosedLoopReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	verbosityLog := new(mytools.VerbosityLog)
	l := verbosityLog.FromContext(ctx)
	l.SetMaxLevel(2)
	l.V(1).Info("Enter Reconcile ClosedLoop")

	//Retreiving ClosedLoop Object who triggered the Reconciler
	closedLoop := &closedlooppoocv1.ClosedLoop{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, closedLoop)

	// Update Status if it's not already done
	// Here we don't use Status for our logic, so the section is not needed but only to show how to use it
	updateStatus := false
	if closedLoop.Name != closedLoop.Status.Name {
		closedLoop.Status.Name = closedLoop.Name
		updateStatus = true
	}

	if closedLoop.Status.IncreaseRank == "" {
		l.V(2).Info("Set start")
		closedLoop.Status.IncreaseRank = "start"
		closedLoop.Status.IncreaseTime = time.Now().String()
		updateStatus = true
	}

	if updateStatus {
		if err := r.Status().Update(ctx, closedLoop); err != nil {
			l.Error(err, "Failed to update closedLoop status")
			return ctrl.Result{}, err
		}
		l.V(2).Info("Enter Reconcile View Spec & Status", "spec", closedLoop.Spec, "status", closedLoop.Status)
	}

	// If the closedLoop who triggered is not find it means that it's been deleted
	if err != nil {

		// The Closedloop object has been deleted or is not found, so we should delete the associated CR (if they exist)
		l.V(2).Info("Close loop Instance not found, Deletion Close Loop ressources")

		if err := r.deleteClosedLoopMonitoringType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No Monitoring found, no deletions have been made")

		}
		if err := r.deleteClosedLoopMonitoringType2(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No Monitoringv2 found, no deletions have been made")

		}

		if err := r.deleteClosedLoopDecisionType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No Decision type 1 found, no deletions have been made")

		}

		if err := r.deleteClosedLoopExecutionType1(ctx, req.Name, req.Namespace, l); err != nil {
			l.Error(err, "No Execution type 1 found, no deletions have been made")

		}

		return ctrl.Result{}, nil

	}
	updateSpec := false
	betterrank := closedLoop.Status.IncreaseRank
	l.V(1).Info("Closedloop receive message: " + "increase rank " + betterrank)

	if updateSpec {
		if err := r.Update(ctx, closedLoop); err != nil {
			l.Error(err, "Failed to update closedLoop spec")
			return ctrl.Result{}, err
		}
		l.V(2).Info("Enter Reconcile View Spec & Status", "spec", closedLoop.Spec, "status", closedLoop.Status)
	}

	// Creation of CR implementation layer ressources (Second layer) are created based on Kinds defined in the primary closedloop ressource
	// These functions are called each time, so if they've already been created, nothing will happen.
	s := mytools.MyReflectValue{Value: reflect.ValueOf(&closedLoop.Spec).Elem()}
	typeOfT := s.Type()
	l.V(1).Info(typeOfT.Kind().String())

	if typeOfT.Kind().String() == "struct" {
		l.V(1).Info(strconv.Itoa(s.NumField()))
		for i := 0; i < s.NumField(); i++ {
			l.V(1).Info(s.Type().Field(i).Name)
			nameCR := s.Type().Field(i).Name
			kindCR := "" //"Monitoringv2"
			//l.V(1).Info("info " + s.XFieldByName(nameCR).XFieldByName("Kind").Kind().String())
			if s.XFieldByName(nameCR).XFieldByName("Kind").Kind() != 0 {

				kindCR = s.XFieldByName(nameCR).XFieldByName("Kind").String()
			}
			l.V(1).Info("kindCR " + kindCR)
			if err := r.createCR(false, nameCR, kindCR, ctx, closedLoop, l); err != nil {
				l.Error(err, "Failed to createCR "+typeOfT.Field(i).Name)
			}
		}
	}
	return ctrl.Result{}, nil
}

// Function use to create the CR
func (r *ClosedLoopReconciler) createCR(ex bool, nameCR string, kindCR string, ctx context.Context, CL *closedlooppoocv1.ClosedLoop, l mytools.VerbosityLog) error {
	if ex {
		return nil
	}
	// Try to retrieve the CR to see if it's already within the Cluster
	clcm := mytools.Cr_type(nameCR, kindCR)
	if clcm == nil {
		return nil
	}
	err := r.Get(ctx, types.NamespacedName{Name: CL.Name + "-" + strings.ToLower(nameCR), Namespace: CL.Namespace}, clcm)

	if err == nil {
		l.V(1).Info("CR Found - No Creation")

		v := mytools.MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}
		cl := reflect.ValueOf(CL).Elem()
		if cl.FieldByName("Spec").FieldByName(nameCR).FieldByName("Kind").Kind() != 0 &&
			cl.FieldByName("Spec").FieldByName(nameCR).FieldByName("Policy").Kind() != 0 &&
			v.FieldByName("Spec").FieldByName("Policy").Kind() != 0 {
			v.FieldByName("Spec").FieldByName("Policy").Set(cl.FieldByName("Spec").FieldByName(nameCR).FieldByName("Policy"))
		}
		if cl.FieldByName("Spec").FieldByName(nameCR).FieldByName("Kind").Kind() != 0 &&
			cl.FieldByName("Spec").FieldByName(nameCR).FieldByName("Config").Kind() != 0 &&
			v.FieldByName("Spec").FieldByName("Config").Kind() != 0 {
			v.FieldByName("Spec").FieldByName("Config").Set(cl.FieldByName("Spec").FieldByName(nameCR).FieldByName("Config"))
		}

		if err := r.Update(ctx, clcm); err != nil {
			l.Error(err, "Failed to update policy in CR "+nameCR)
		}
		return nil
	}
	l.V(1).Info("CR not found, Creating CR " + nameCR)

	v := mytools.MyReflectValue{Value: reflect.ValueOf(clcm).Elem()}

	name := CL.Name + "-" + strings.ToLower(nameCR)
	m := "{}"

	l.V(1).Info("Name " + name)
	l.V(1).Info("Namespace " + v.FieldByName("Namespace").String())
	v.FieldByName("Namespace").Set(reflect.ValueOf(&CL.Namespace).Elem())
	v.FieldByName("Name").Set(reflect.ValueOf(&name).Elem())
	v.FieldByName("Spec").FieldByName("Affix").Set(reflect.ValueOf(&CL.Name).Elem())
	l.V(1).Info("Time " + v.FieldByName("Spec").FieldByName("Time").String())
	v.XFieldByName("Spec").XFieldByName("Time").XSetString(time.Now().String())
	v.FieldByName("Spec").FieldByName("Message").Set(reflect.ValueOf(&m).Elem())

	return r.Create(ctx, clcm)

}

func (r *ClosedLoopReconciler) deleteClosedLoopMonitoringType1(ctx context.Context, name string, namespace string, l mytools.VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-monitoring"
	clcm := &closedlooppoocv1.Monitoring{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.V(2).Info("Monitoring type 1 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete monitoring")
		return err
	}

	l.V(2).Info("Deleted monitoring")

	return nil

}

func (r *ClosedLoopReconciler) deleteClosedLoopMonitoringType2(ctx context.Context, name string, namespace string, l mytools.VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-monitoring"
	clcm := &closedlooppoocv1.Monitoringv2{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.V(2).Info("Monitoring type 2 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete monitoring")
		return err
	}

	l.Info("Deleted monitoring")

	return nil

}

func (r *ClosedLoopReconciler) deleteClosedLoopDecisionType1(ctx context.Context, name string, namespace string, l mytools.VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-decision"
	clcm := &closedlooppoocv1.Decision{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.V(2).Info("Decision type 1 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete decision")
		return err
	}

	l.V(2).Info("Deleted decision")

	return nil

}

func (r *ClosedLoopReconciler) deleteClosedLoopExecutionType1(ctx context.Context, name string, namespace string, l mytools.VerbosityLog) error {
	// Try to retrieve the CR to see if it's already within the Cluster
	name = name + "-execution"
	clcm := &closedlooppoocv1.Execution{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, clcm)

	if err != nil {
		if errors.IsNotFound(err) {
			l.V(2).Info("Execution type 1 already deleted or no present")
			return nil
		}
		return err
	}
	//Delete Object if present within the CLuster
	err = r.Delete(ctx, clcm)
	if err != nil {
		l.Error(err, "Failed to delete execution")
		return err
	}
	l.V(2).Info("Deleted execution")

	return nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *ClosedLoopReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&closedlooppoocv1.ClosedLoop{}).
		Complete(r)
}

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

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	schedulev1 "github.com/techierishi/k8soperator/api/v1"
)

// CrudReconciler reconciles a Crud object
type CrudReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=schedule.rs,resources=cruds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=schedule.rs,resources=cruds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=schedule.rs,resources=cruds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Crud object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *CrudReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("Crud", req.NamespacedName)

	// Fetch the Crud instance
	instance := &schedulev1.Crud{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
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

	pvFound := &corev1.PersistentVolume{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.App.Name, Namespace: instance.Namespace}, pvFound)
	var pvResult *reconcile.Result
	pvResult, err = r.ensurePersistentVolume(req, instance, r.persistentVolume(instance.Spec.Volume, instance))
	if pvResult != nil {
		log.Error(err, "Persistent Volume Not ready")
		return *pvResult, err
	}

	// PV already exists - don't requeue
	log.Info("Skip reconcile: PV already exists",
		"PV.Namespace", pvFound.Namespace, "PV.Name", pvFound.Name)
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.App.Name, Namespace: instance.Namespace}, pvcFound)
	var pvcResult *reconcile.Result
	pvcResult, err = r.ensurePersistentVolumeClaim(req, instance, r.persistentVolumeClaim(instance.Spec.Volume, instance))
	if pvcResult != nil {
		log.Error(err, "Persistent Volume Not ready")
		return *pvcResult, err
	}
	// PVC already exists - don't requeue
	log.Info("Skip reconcile: PVC already exists",
		"PVC.Namespace", pvcFound.Namespace, "PVC.Name", pvcFound.Name)

	// Check if this Deployment already exists
	appFound := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.App.Name, Namespace: instance.Namespace}, appFound)
	var appResult *reconcile.Result
	appResult, err = r.ensureDeployment(req, instance, r.backendDeployment(instance.Spec.App, nil, instance))
	if appResult != nil {
		log.Error(err, "App Deployment Not ready")
		return *appResult, err
	}

	// Check if this Service already exists
	appResult, err = r.ensureService(req, instance, r.backendService(instance.Spec.App, instance))
	if appResult != nil {
		log.Error(err, "App Service Not ready")
		return *appResult, err
	}
	// App Deployment and Service already exists - don't requeue
	log.Info("Skip reconcile: App Deployment and service already exists",
		"Deployment.Namespace", appFound.Namespace, "Deployment.Name", appFound.Name)

	dbFound := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Db.Name, Namespace: instance.Namespace}, dbFound)
	var dbResult *reconcile.Result
	dbResult, err = r.ensureDeployment(req, instance, r.backendDeployment(instance.Spec.Db, &instance.Spec.Volume, instance))
	if dbResult != nil {
		log.Error(err, "Db Deployment Not ready")
		return *dbResult, err
	}

	dbResult, err = r.ensureService(req, instance, r.backendService(instance.Spec.Db, instance))
	if dbResult != nil {
		log.Error(err, "Db Service Not ready")
		return *dbResult, err
	}

	// Db Deployment and Service already exists - don't requeue
	log.Info("Skip reconcile: DB Deployment and service already exists",
		"Deployment.Namespace", dbFound.Namespace, "Deployment.Name", dbFound.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CrudReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&schedulev1.Crud{}).
		Complete(r)
}

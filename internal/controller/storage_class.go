package controller

import (
	"context"

	mydomainv1alpha1 "github.com/techierishi/k8soperator/api/v1"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensureSc ensures Sc is Running in a namespace.
func (r *CrudReconciler) ensureStorageClass(request reconcile.Request,
	instance *mydomainv1alpha1.Crud,
	sc *storagev1.StorageClass,
) (*reconcile.Result, error) {

	// See if sc already exists and create if it doesn't
	found := &storagev1.StorageClass{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      sc.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the sc
		err = r.Create(context.TODO(), sc)

		if err != nil {
			// Sc creation failed
			return &reconcile.Result{}, err
		} else {
			// Sc creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the sc not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendSc is a code for creating a Sc
func (r *CrudReconciler) storageClass(vol mydomainv1alpha1.Volume, v *mydomainv1alpha1.Crud) *storagev1.StorageClass {

	volumeExpansion := true
	volumeBindingMode := storagev1.VolumeBindingMode("WaitForFirstConsumer")
	reclaimPolicy := corev1.PersistentVolumeReclaimPolicy("Retain")
	sc := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "crud-storage-class",
			Namespace: v.Namespace,
		},
		VolumeBindingMode:    &volumeBindingMode,
		AllowVolumeExpansion: &volumeExpansion,
		ReclaimPolicy:        &reclaimPolicy,
		Provisioner:          "kubernetes.io/no-provisioner",
	}

	controllerutil.SetControllerReference(v, sc, r.Scheme)
	return sc
}

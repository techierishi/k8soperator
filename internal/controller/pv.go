package controller

import (
	"context"

	mydomainv1alpha1 "github.com/techierishi/k8soperator/api/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensurePv ensures Pv is Running in a namespace.
func (r *CrudReconciler) ensurePersistentVolume(request reconcile.Request,
	instance *mydomainv1alpha1.Crud,
	pv *corev1.PersistentVolume,
) (*reconcile.Result, error) {

	// See if pv already exists and create if it doesn't
	found := &corev1.PersistentVolume{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      pv.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the pv
		err = r.Create(context.TODO(), pv)

		if err != nil {
			// Pv creation failed
			return &reconcile.Result{}, err
		} else {
			// Pv creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the pv not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendPv is a code for creating a Pv
func (r *CrudReconciler) persistentVolume(vol mydomainv1alpha1.Volume, v *mydomainv1alpha1.Crud) *corev1.PersistentVolume {

	volumeMode := corev1.PersistentVolumeMode("Filesystem")
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vol.PvName,
			Namespace: v.Namespace,
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(vol.Capacity),
			},
			VolumeMode: &volumeMode,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: vol.Path,
				},
			},
			StorageClassName: "crud-storage-class",
		},
	}

	controllerutil.SetControllerReference(v, pv, r.Scheme)
	return pv
}

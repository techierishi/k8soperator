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

// ensurePvc ensures Pvc is Running in a namespace.
func (r *CrudReconciler) ensurePersistentVolumeClaim(request reconcile.Request,
	instance *mydomainv1alpha1.Crud,
	pvc *corev1.PersistentVolumeClaim,
) (*reconcile.Result, error) {

	// See if pvc already exists and create if it doesn't
	found := &corev1.PersistentVolumeClaim{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      pvc.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the pvc
		err = r.Create(context.TODO(), pvc)

		if err != nil {
			// Pvc creation failed
			return &reconcile.Result{}, err
		} else {
			// Pvc creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the pvc not existing
		return &reconcile.Result{}, err
	}

	if found.Spec.Resources.Requests.Storage() != pvc.Spec.Resources.Requests.Storage() {
		// Create the pvc
		err = r.Create(context.TODO(), pvc)

		if err != nil {
			// Pvc creation failed
			return &reconcile.Result{}, err
		} else {
			// Pvc creation was successful
			return nil, nil
		}
	}

	return nil, nil
}

// backendPvc is a code for creating a Pvc
func (r *CrudReconciler) persistentVolumeClaim(vol mydomainv1alpha1.Volume, v *mydomainv1alpha1.Crud) *corev1.PersistentVolumeClaim {

	storageCls := "crud-storage-class"
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vol.PvcName,
			Namespace: v.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(vol.Storage),
				},
			},
			StorageClassName: &storageCls,
		},
	}

	controllerutil.SetControllerReference(v, pvc, r.Scheme)
	return pvc
}

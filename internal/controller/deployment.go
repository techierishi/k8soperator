package controller

import (
	"context"

	mydomainv1alpha1 "github.com/techierishi/k8soperator/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func labels(v *mydomainv1alpha1.Crud, tier string) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app":             "visitors",
		"visitorssite_cr": v.Name,
		"tier":            tier,
	}
}

// ensureDeployment ensures Deployment resource presence in given namespace.
func (r *CrudReconciler) ensureDeployment(request reconcile.Request,
	instance *mydomainv1alpha1.Crud,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendDeployment is a code for Creating Deployment
func (r *CrudReconciler) backendDeployment(module mydomainv1alpha1.Module, vol *mydomainv1alpha1.Volume, v *mydomainv1alpha1.Crud) *appsv1.Deployment {

	volMounts := []corev1.VolumeMount{}
	volumes := []corev1.Volume{}

	if vol != nil {

		volMounts = []corev1.VolumeMount{
			{
				MountPath: vol.Path,
				Name:      vol.PvcName,
			},
		}

		volumes = []corev1.Volume{
			{
				Name: vol.PvcName,
			},
		}
	}

	labels := labels(v, module.Name)
	size := int32(1)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      module.Name,
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           module.Image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            module.Name,
						Ports: []corev1.ContainerPort{{
							ContainerPort: module.Port,
						}},
						VolumeMounts: volMounts,
					}},

					Volumes: volumes,
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

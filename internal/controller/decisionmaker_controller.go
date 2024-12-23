/*
Copyright 2024.

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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mecdmsv1alpha1 "github.com/adeptvin1/lwmecps-operator/api/v1alpha1"
)

// DecisionMakerReconciler reconciles a DecisionMaker object
type DecisionMakerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mecdms.apps.lwmecps.com,resources=decisionmakers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mecdms.apps.lwmecps.com,resources=decisionmakers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mecdms.apps.lwmecps.com,resources=decisionmakers/finalizers,verbs=update

// Reconcile reconciles the DecisionMaker resource
func (r *DecisionMakerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Получение ресурса DecisionMaker
	var decisionMaker mecdmsv1alpha1.DecisionMaker
	if err := r.Get(ctx, req.NamespacedName, &decisionMaker); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("DecisionMaker resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get DecisionMaker")
		return ctrl.Result{}, err
	}

	// Определяем количество реплик
	replicas := int32(1)
	if decisionMaker.Spec.Replicas != nil {
		replicas = *decisionMaker.Spec.Replicas
	}

	// Создание Deployment для Nginx
	mecDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      decisionMaker.Name + "-mec",
			Namespace: req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": decisionMaker.Name + "-mec"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": decisionMaker.Name + "-mec"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mec",
							Image: "nginx",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Проверяем, существует ли Deployment
	var existingDeployment appsv1.Deployment
	err := r.Get(ctx, types.NamespacedName{Name: mecDeployment.Name, Namespace: mecDeployment.Namespace}, &existingDeployment)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("Creating a new Deployment", "Deployment.Namespace", mecDeployment.Namespace, "Deployment.Name", mecDeployment.Name)
			if err := r.Create(ctx, mecDeployment); err != nil {
				logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", mecDeployment.Namespace, "Deployment.Name", mecDeployment.Name)
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		logger.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Если Deployment уже существует
	logger.Info("Deployment already exists", "Deployment.Namespace", mecDeployment.Namespace, "Deployment.Name", mecDeployment.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DecisionMakerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mecdmsv1alpha1.DecisionMaker{}).
		Complete(r)
}

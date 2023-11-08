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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mydomainv1 "github.com/ImmortalHedgehogs/TamaKubei/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestReconciler reconciles a Test object
type TestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.domain,resources=tests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.domain,resources=tests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.domain,resources=tests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Test object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
// func (r *TestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
// 	log := log.FromContext(ctx)

// 	var test mydomainv1.Test
// 	if err := r.Get(ctx, req.NamespacedName, &test); err != nil {
// 		if client.IgnoreNotFound(err) == nil {
// 			// Resource has been deleted. Delete the corresponding Job.
// 			// if err := r.deleteJob(ctx, test); err != nil {
// 			// 	log.Error(err, "unable to delete Job")
// 			// 	return ctrl.Result{}, err
// 			// }
// 			jobSpec, err := r.deleteJob(ctx, test)
// 			if err != nil {
// 				log.Error(err, "failed to delete Job spec")
// 				return ctrl.Result{}, client.IgnoreNotFound(err)
// 			}
// 			if err := r.Delete(ctx, &jobSpec); err != nil {
// 				log.Error(err, "unable to delete Job")
// 			}
// 		} else {
// 			log.Error(err, "unable to fetch Test")
// 		}
// 		return ctrl.Result{}, client.IgnoreNotFound(err)
// 	}

// 	jobSpec, err := r.createJob(test)
// 	if err != nil {
// 		log.Error(err, "failed to create Job spec")
// 		return ctrl.Result{}, client.IgnoreNotFound(err)
// 	}

// 	if err := r.Create(ctx, &jobSpec); err != nil {
// 		log.Error(err, "unable to create Job")
// 	}

// 	return ctrl.Result{}, nil
// }

func (r *TestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var test mydomainv1.Test
	if err := r.Get(ctx, req.NamespacedName, &test); err != nil {
		log.Error(err, "unable to fetch Test")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	jobSpec, err := r.createJob(test)
	if err != nil {
		log.Error(err, "failed to create Job spec")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.Create(ctx, &jobSpec); err != nil {
		log.Error(err, "unable to create Job")
	}

	return ctrl.Result{}, nil
}

func (r *TestReconciler) createJob(test mydomainv1.Test) (batchv1.Job, error) {
	j := batchv1.Job{
		TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "Test"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.Name + "-job",
			Namespace: test.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:    "main-container",
							Image:   "alpine",
							Command: []string{"sh", "-c", "sleep 3600"},
						},
					},
				},
			},
		},
	}

	return j, nil
}

func (r *TestReconciler) deleteJob(ctx context.Context, test mydomainv1.Test) (batchv1.Job, error) {
	jobName := test.Name + "-job"
	jobNamespace := test.Namespace

	// Create a Job key
	jobKey := client.ObjectKey{
		Name:      jobName,
		Namespace: jobNamespace,
	}

	// Check if the Job exists before attempting to delete it
	var existingJob batchv1.Job
	if err := r.Get(ctx, jobKey, &existingJob); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// The Job exists, so delete it
			return existingJob, nil
		} else {
			// The Job does not exist, or there was an error
			return existingJob, err
		}
	}

	return existingJob, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mydomainv1.Test{}).
		Complete(r)
}

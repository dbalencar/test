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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testv1alpha1 "github.com/dbalencar/test-operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ExternalReconciler reconciles a External object
type ExternalReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const externalFinalizer = "test.dbalencar.com/finalizer"

//+kubebuilder:rbac:groups=test.dbalencar.com,resources=externals,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.dbalencar.com,resources=externals/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.dbalencar.com,resources=externals/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the External object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ExternalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx).WithValues("external", req.NamespacedName)
	reqLogger.Info("Reconciling External")

	// TODO(user): your logic here
	external := &testv1alpha1.External{}
	err := r.Get(ctx, req.NamespacedName, external)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("External resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get External.")
		return ctrl.Result{}, err
	}

	isExternalMarkedToBeDeleted := external.GetDeletionTimestamp() != nil
	if isExternalMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(external, externalFinalizer) {

			if err := r.finalizeExternal(reqLogger, external); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(external, externalFinalizer)
			err := r.Update(ctx, external)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(external, externalFinalizer) {
		controllerutil.AddFinalizer(external, externalFinalizer)
		err = r.Update(ctx, external)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1alpha1.External{}).
		Complete(r)
}

func (r *ExternalReconciler) finalizeExternal(reqLogger logr.Logger, e *testv1alpha1.External) error {

	reqLogger.Info("Successfully finalized external")
	return nil
}

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
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appv1 "kubebuilder.io/apps/api/v1"
	//coreV1 "k8s.io/api/core/v1"
)

// DeployObjectReconciler reconciles a DeployObject object
type DeployObjectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.kubebuilder.io,resources=deployobjects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.kubebuilder.io,resources=deployobjects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.kubebuilder.io,resources=deployobjects/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeployObject object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile

func EventFunc(r *DeployObjectReconciler, req *appv1.DeployObject) error {

	e := ObjectCheck(req)
	if e == nil {
		volumes := CreateVolumeForDeployment(req)
		configMap := CreateConfigMap(req)
		deployment := CreateDeployment(req, volumes)
		service := CreateService(req)
		k8sControll := NewK8s_resource_controller(r, req, configMap, deployment, service)

		if req.DeletionTimestamp != nil {
			fmt.Println("del ----:")
			k8sControll.Delete()
		} else if req.DeletionTimestamp == nil && req.Status.Status == 0 {
			fmt.Println("add or update---:")
			k8sControll.AddOrUpdate()
		} else if req.Status.Status == 1 {
			return errors.New("无需处理")
		}
	} else {
		req.Status.Describe = e.Error()
		req.Status.Status = 1
		req.Finalizers = req.Finalizers[:0]
		r.Update(context.Background(), req)
	}
	return nil
}

func (r *DeployObjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	var deployObject appv1.DeployObject

	error := r.Get(ctx, req.NamespacedName, &deployObject)

	if error != nil {
		fmt.Println("error----:", error)
		return ctrl.Result{}, nil
	}

	EventFunc(r, &deployObject)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeployObjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	p := predicate.Funcs{
		CreateFunc: func(ce event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(de event.DeleteEvent) bool {
			return true
		},
		UpdateFunc: func(ue event.UpdateEvent) bool {
			return true
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.DeployObject{}).WithEventFilter(p).
		Complete(r)
}

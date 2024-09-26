package controller

import (
	"context"
	"errors"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"

	//"k8s.io/apimachinery/pkg/api/meta"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	appv1 "kubebuilder.io/apps/api/v1"
	//"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8s_resource_controller struct {
	R          *DeployObjectReconciler
	Object     *appv1.DeployObject
	ConfigMap  *coreV1.ConfigMap
	Deployment *appsV1.Deployment
	Service    *coreV1.Service
	Time       context.Context
}

func NewK8s_resource_controller(r *DeployObjectReconciler, o *appv1.DeployObject, c *coreV1.ConfigMap, d *appsV1.Deployment, s *coreV1.Service) *K8s_resource_controller {

	return &K8s_resource_controller{
		R:          r,
		Object:     o,
		ConfigMap:  c,
		Deployment: d,
		Service:    s,
		Time:       context.Background(),
	}

}

func (k8s *K8s_resource_controller) Delete() {
	// 再删除crd
	k8s.Object.Finalizers = k8s.Object.Finalizers[:0]
	k8s.R.Update(context.Background(), k8s.Object)

	// 先删除crd 定义的资源
	k8s.R.Delete(k8s.Time, k8s.Deployment)
	k8s.R.Delete(k8s.Time, k8s.ConfigMap)
	k8s.R.Delete(k8s.Time, k8s.Service)

}

func (k8s *K8s_resource_controller) AddOrUpdate() {

	// 创建或者更新configmap
	x := types.NamespacedName{
		Namespace: k8s.Object.Namespace,
		Name:      k8s.Object.Name,
	}
	configMap := &coreV1.ConfigMap{}
	e := k8s.R.Get(k8s.Time, x, configMap)
	if e != nil {
		e := k8s.R.Create(k8s.Time, k8s.ConfigMap)
		if e != nil {
			k8s.Object.Status.Describe = "ERROR(create configmap) :" + e.Error()
			k8s.Object.Status.Status = 1
			k8s.R.Update(k8s.Time, k8s.Object)
			return
		} else {
			k8s.Object.Status.Describe = "SUCCESS:create configmap success"
			k8s.Object.Status.Status = 1
			k8s.R.Update(k8s.Time, k8s.Object)
		}
	} else {
		e := k8s.R.Update(k8s.Time, k8s.ConfigMap)
		if e != nil {
			k8s.Object.Status.Describe = "ERROR(update configmap) :" + e.Error()
			k8s.Object.Status.Status = 1
			k8s.R.Update(k8s.Time, k8s.Object)
			return
		} else {
			k8s.Object.Status.Describe = "SUCCESS:update configMap success"
			k8s.Object.Status.Status = 1
			k8s.R.Update(k8s.Time, k8s.Object)
		}
	}

	//创建或者更新deployment
	deployment := &appsV1.Deployment{}
	er := k8s.R.Get(context.Background(), x, deployment)
	if er != nil {
		e := k8s.R.Create(k8s.Time, k8s.Deployment)
		if e != nil {
			k8s.Object.Status.Describe = "ERROR(create deployment) :" + e.Error()
			k8s.Object.Status.Status = 1
			k8s.R.Update(k8s.Time, k8s.Object)
			return
		} else {
			k8s.Object.Status.Describe = "SUCCESS:create deployment success"
			k8s.Object.Status.Status = 1
			k8s.R.Update(k8s.Time, k8s.Object)
		}
	} else {
		e := k8s.R.Update(k8s.Time, k8s.Deployment)
		if e != nil {
			k8s.Object.Status.Describe = "ERROR(update deployment) :" + e.Error()
			k8s.Object.Status.Status = 1
			k8s.R.Update(k8s.Time, k8s.Object)
			return
		} else {
			k8s.Object.Status.Describe = "SUCCESS:update deployment success"
			k8s.Object.Status.Status = 1
			k8s.R.Update(k8s.Time, k8s.Object)
		}
	}

	//创建或者更新service
	service := &coreV1.Service{}
	e = k8s.R.Get(k8s.Time, x, service)
	if e != nil {
		e = k8s.R.Create(k8s.Time, k8s.Service)
		k8s.Object.Status.Describe = "SUCCESS:create service success"
		k8s.Object.Status.Status = 1
		k8s.R.Update(k8s.Time, k8s.Object)
	} else {
		e = k8s.R.Update(k8s.Time, k8s.Service)
		k8s.Object.Status.Describe = "SUCCESS:update service success"
		k8s.Object.Status.Status = 1
		k8s.R.Update(k8s.Time, k8s.Object)
	}

}

func ObjectCheck(object *appv1.DeployObject) error {
	if len(object.Spec.Image) == 0 {
		return errors.New("缺少spect.image 参数")
	}
	if len(object.Spec.AppLogDir) == 0 {
		return errors.New("缺少spect.appLogDir 参数")
	}
	if object.Spec.Port == 0 {
		return errors.New("缺少spect.port 参数")
	}

	if object.Spec.Num == 0 {
		object.Spec.Num = 1
	}

	if len(object.Spec.Labels) == 0 {
		object.Spec.Labels = make(map[string]string)
		object.Spec.Labels["app"] = object.Name
	}

	if len(object.Spec.Cpu) == 0 {
		object.Spec.Cpu = "500m"
	}
	if len(object.Spec.Mem) == 0 {
		object.Spec.Mem = "1Gi"
	}

	if len(object.Finalizers) == 0 {
		object.Finalizers = []string{"x"}
	}

	if len(object.Spec.LogFormat) == 0 {
		object.Spec.LogFormat = `^\[`
	}
    
	if len(object.Spec.ResourceLevel) == 0 {
		object.Spec.ResourceLevel = "0"
	}
	return nil

}

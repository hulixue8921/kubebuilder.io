package controller

import (
	coreV1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	appv1 "kubebuilder.io/apps/api/v1"
)

func CreateService(object *appv1.DeployObject) *coreV1.Service {

	req := object.DeepCopy()

	return &coreV1.Service{
		TypeMeta: meta.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: coreV1.ServiceSpec{
			Selector: req.Spec.Labels,
			Ports: []coreV1.ServicePort{
				coreV1.ServicePort{
					Port:       req.Spec.Port,
					Protocol:   coreV1.ProtocolTCP,
					TargetPort: intstr.FromInt(int(req.Spec.Port)),
				},
			},
		},
	}

}

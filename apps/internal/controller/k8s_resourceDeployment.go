package controller

import (
	"time"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	appv1 "kubebuilder.io/apps/api/v1"
)

const layout = "20060102150405"

func CreateLogContain(object *appv1.DeployObject) *coreV1.Container {
	req := object.DeepCopy()
	cpuQuantity, _ := resource.ParseQuantity("500m")
	memQuantity, _ := resource.ParseQuantity("1Gi")
	r := coreV1.ResourceList{
		coreV1.ResourceCPU:    cpuQuantity,
		coreV1.ResourceMemory: memQuantity,
	}

	return &coreV1.Container{
		Name:            req.Name + "-log",
		Image:           "logstash:7.4.1",
		ImagePullPolicy: "IfNotPresent",
		Resources: coreV1.ResourceRequirements{
			Limits:   r,
			Requests: r,
		},
		Lifecycle: &coreV1.Lifecycle{
			PreStop: &coreV1.LifecycleHandler{
				Exec: &coreV1.ExecAction{
					Command: []string{
						"sleep 10",
					},
				},
			},
		},
		VolumeMounts: []coreV1.VolumeMount{
			coreV1.VolumeMount{
				Name:      "time",
				MountPath: "/etc/localtime",
			},
			coreV1.VolumeMount{
				Name:      "sharelog",
				MountPath: req.Spec.AppLogDir,
			},
			coreV1.VolumeMount{
				Name:      "logstashyml",
				MountPath: "/usr/share/logstash/config",
			},
			coreV1.VolumeMount{
				Name:      "logconf",
				MountPath: "/usr/share/logstash/conf.d",
			},
		},
	}

}

func CreateContain(object *appv1.DeployObject) *coreV1.Container {
	req := object.DeepCopy()
	// 资源要求
	cpuQuantity, _ := resource.ParseQuantity(req.Spec.Cpu)
	memQuantity, _ := resource.ParseQuantity(req.Spec.Mem)
	r := coreV1.ResourceList{
		coreV1.ResourceCPU:    cpuQuantity,
		coreV1.ResourceMemory: memQuantity,
	}
     
	var resource coreV1.ResourceRequirements
	if req.Spec.ResourceLevel == "0" {
        resource =coreV1.ResourceRequirements{
              Limits: r,
		}
	} else {
		resource =coreV1.ResourceRequirements{
			Limits: r,
			Requests: r,
		}
	}

	// 需要挂载的卷
	volumeMounts := []coreV1.VolumeMount{}
	timeMounts := coreV1.VolumeMount{
		Name:      "time",
		MountPath: "/etc/localtime",
	}

	logMounts := coreV1.VolumeMount{
		Name:      "sharelog",
		MountPath: req.Spec.AppLogDir,
	}

	if len(req.Spec.Disk.Path) == 0 {
		volumeMounts = append(volumeMounts, timeMounts, logMounts)
	} else {
		otherMounts := coreV1.VolumeMount{
			Name:      req.Name + "-nfs",
			MountPath: req.Spec.Disk.Path,
		}
		volumeMounts = append(volumeMounts, timeMounts, logMounts, otherMounts)
	}

	if len(req.Spec.Secret) != 0 {
		for _, name := range req.Spec.Secret {

			secretMount := coreV1.VolumeMount{
				Name:      name,
				MountPath: "/etc/secret/" + name,
			}
			volumeMounts = append(volumeMounts, secretMount)
		}
	}

	// 探针对象
	probe := &coreV1.Probe{
		InitialDelaySeconds: 15,
		PeriodSeconds:       20,
		ProbeHandler: coreV1.ProbeHandler{
			TCPSocket: &coreV1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: req.Spec.Port,
				},
			},
		},
	}
	return &coreV1.Container{
		Name:            req.Name,
		Image:           req.Spec.Image,
		ImagePullPolicy: "IfNotPresent",
		Env: []coreV1.EnvVar{
			coreV1.EnvVar{
				Name:  "restart",
				Value: time.Now().Format(layout),
			},
		},
		Ports: []coreV1.ContainerPort{
			coreV1.ContainerPort{
				ContainerPort: req.Spec.Port,
				Protocol:      coreV1.ProtocolTCP,
			},
		},
		LivenessProbe:  probe,
		ReadinessProbe: probe,
		Lifecycle: &coreV1.Lifecycle{
			PostStart: &coreV1.LifecycleHandler{
				Exec: &coreV1.ExecAction{
					Command: []string{
						"./start.sh",
					},
				},
			},
			PreStop: &coreV1.LifecycleHandler{
				Exec: &coreV1.ExecAction{
					Command: []string{
						"./stop.sh",
					},
				},
			},
		},
		/*
		Resources: coreV1.ResourceRequirements{
			Limits:   r,
			Requests: r,
		},
		*/
	    Resources: resource,
		VolumeMounts: volumeMounts,
	}

}

func CreateDeployment(object *appv1.DeployObject, vs []VolumeInterface) *appsV1.Deployment {
	pullSecret := "regcred"
	imagePullSecrets := []coreV1.LocalObjectReference{
		coreV1.LocalObjectReference{
			Name: pullSecret,
		},
	}
	req := object.DeepCopy()
	//req.Spec.Labels["app"] = req.Name

	appPod := CreateContain(req)
	logPod := CreateLogContain(req)

	volumes := []coreV1.Volume{}

	for _, v := range vs {
		volumes = append(volumes, *v.Create())
	}
	return &appsV1.Deployment{
		TypeMeta: meta.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Spec.Labels,
			Annotations: req.Spec.Annotations,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: &req.Spec.Num,
			Selector: &meta.LabelSelector{
				MatchLabels: req.Spec.Labels,
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: req.Spec.Labels,
				},
				Spec: coreV1.PodSpec{
					ImagePullSecrets: imagePullSecrets,
					Affinity: &coreV1.Affinity{
						PodAntiAffinity: &coreV1.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []coreV1.PodAffinityTerm{
								coreV1.PodAffinityTerm{LabelSelector: &meta.LabelSelector{
									MatchExpressions: []meta.LabelSelectorRequirement{
										meta.LabelSelectorRequirement{
											Key:      "app",
											Operator: "In",
											Values: []string{
												req.Name,
											},
										},
									},
								},
									TopologyKey: "kubernetes.io/hostname",
								},
							},
						},
					},
					Containers: []coreV1.Container{
						*appPod,
						*logPod,
					},
					Volumes: volumes,
				},
			},
		},
	}

}

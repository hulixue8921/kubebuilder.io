package controller

import (
	coreV1 "k8s.io/api/core/v1"
	appv1 "kubebuilder.io/apps/api/v1"
)

const (
	NfsServer   = "192.168.1.1"
	NfsRootPath = "/data/nfs/"
)

type VolumeInterface interface {
	Create() *coreV1.Volume
}

type HostV struct {
	Name string
	Path string
}

func (v *HostV) Create() *coreV1.Volume {
	return &coreV1.Volume{
		Name: v.Name,
		VolumeSource: coreV1.VolumeSource{
			HostPath: &coreV1.HostPathVolumeSource{
				Path: v.Path,
			},
		},
	}
}

type ShareV struct {
	Name string
}

func (v *ShareV) Create() *coreV1.Volume {
	return &coreV1.Volume{
		Name: v.Name,
		VolumeSource: coreV1.VolumeSource{
			EmptyDir: &coreV1.EmptyDirVolumeSource{},
		},
	}

}

type ConfigMapV struct {
	Name          string
	ConfigMapName string
	ConfigMapKey  string
	Path          string
}

func (v *ConfigMapV) Create() *coreV1.Volume {
	return &coreV1.Volume{
		Name: v.Name,
		VolumeSource: coreV1.VolumeSource{
			ConfigMap: &coreV1.ConfigMapVolumeSource{
				LocalObjectReference: coreV1.LocalObjectReference{
					Name: v.ConfigMapName,
				},
				Items: []coreV1.KeyToPath{
					coreV1.KeyToPath{
						Key:  v.ConfigMapKey,
						Path: v.Path,
					},
				},
			},
		},
	}

}

type PvcVolume struct {
	Name    string
	PvcName string
}

func (pvc *PvcVolume) Create() *coreV1.Volume {

	return &coreV1.Volume{
		Name: pvc.Name,
		VolumeSource: coreV1.VolumeSource{
			PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvc.PvcName,
			},
		},
	}
}

type NfsVolume struct {
	Name      string
	NfsServer string
	NfsPath   string
}

func (pvc *NfsVolume) Create() *coreV1.Volume {
	return &coreV1.Volume{
		Name: pvc.Name,
		VolumeSource: coreV1.VolumeSource{
			NFS: &coreV1.NFSVolumeSource{
				Server: pvc.NfsServer,
				Path:   pvc.NfsPath,
			},
		},
	}
}

type SecretVolume struct {
	SecretName string
}

func (secretVolume *SecretVolume) Create() *coreV1.Volume {
	return &coreV1.Volume{
		Name: secretVolume.SecretName,
		VolumeSource: coreV1.VolumeSource{
			Secret: &coreV1.SecretVolumeSource{
				SecretName: secretVolume.SecretName,
			},
		},
	}
}

func CreateVolumeForDeployment(Object *appv1.DeployObject) []VolumeInterface {
	// 需要的卷
	volumes := []VolumeInterface{}
	hostV := &HostV{
		Name: "time",
		Path: "/etc/localtime",
	}
	logV := &ShareV{
		Name: "sharelog",
	}
	configmapV1 := &ConfigMapV{
		Name:          "logstashyml",
		ConfigMapName: Object.Name,
		ConfigMapKey:  "logstash.yml",
		Path:          "logstash.yml",
	}
	configmapV2 := &ConfigMapV{
		Name:          "logconf",
		ConfigMapName: Object.Name,
		ConfigMapKey:  "log.conf",
		Path:          "log.conf",
	}

	if len(Object.Spec.Disk.Path) == 0 {
		volumes = append(volumes, hostV, logV, configmapV1, configmapV2)
	} else {
		nfsVolume := &NfsVolume{
			Name:      Object.Name + "-nfs",
			NfsServer: NfsServer,
			NfsPath:   NfsRootPath + Object.Name,
		}
		volumes = append(volumes, hostV, logV, configmapV1, configmapV2, nfsVolume)
	}

	if len(Object.Spec.Secret) != 0 {
		secretVolume := &SecretVolume{
			SecretName: Object.Spec.Secret,
		}
		volumes = append(volumes, secretVolume)
	}
	return volumes

}

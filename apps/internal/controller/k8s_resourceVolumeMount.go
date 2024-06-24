package controller

import (
	coreV1 "k8s.io/api/core/v1"
)

type VolumeMountInterface interface {
	Mount() *coreV1.VolumeMount
}

type VolumeMount struct {
	VolumeName string
	Path       string
}

func (vM *VolumeMount) Mount() *coreV1.VolumeMount {

	return &coreV1.VolumeMount{
		Name:      vM.VolumeName,
		MountPath: vM.Path,
	}

}

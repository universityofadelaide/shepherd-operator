package image

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Container provides options for a step when creating the image.
type Container struct {
	Image  string
	CPU    string
	Memory string
}

// ResourceRequirements used with this image.
func (c Container) ResourceRequirements() (corev1.ResourceRequirements, error) {
	var requirements corev1.ResourceRequirements

	cpu, err := resource.ParseQuantity(c.CPU)
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}

	memory, err := resource.ParseQuantity(c.Memory)
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}

	requirements.Requests = corev1.ResourceList{
		corev1.ResourceCPU:    cpu,
		corev1.ResourceMemory: memory,
	}

	requirements.Limits = corev1.ResourceList{
		corev1.ResourceCPU:    cpu,
		corev1.ResourceMemory: memory,
	}

	return requirements, nil
}

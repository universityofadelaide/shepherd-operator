package drupal

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
)

// BuildConslidatedConfig from Secrets and ConfigMaps which can be consumed by an application.
func BuildConslidatedConfig(defaultConfig, overrideConfig *corev1.ConfigMap, defaultSecret, overrideSecret *corev1.Secret) ([]byte, error) {
	list := defaultConfig.Data

	for key, value := range overrideConfig.Data {
		list[key] = value
	}

	for key, value := range defaultSecret.Data {
		list[key] = string(value)
	}

	for key, value := range overrideSecret.Data {
		list[key] = string(value)
	}

	return json.Marshal(list)
}

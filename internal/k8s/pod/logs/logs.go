package logs

import (
	"context"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Interface for the logs client.
type Interface interface {
	Get(string, string, string) (string, error)
}

// Client for getting pod logs.
type Client struct {
	kubeset *kubernetes.Clientset
}

// New creates a new client.
func New(config *rest.Config) (Client, error) {
	kubeset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return Client{}, err
	}
	return Client{
		kubeset: kubeset,
	}, nil
}

// Get gets the logs for a container from a pod in a namespace.
func (c Client) Get(ctx context.Context, namespace, name, container string) (string, error) {
	body, err := c.kubeset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{
		Container: container,
	}).Stream(ctx)
	if err != nil {
		return "", err
	}
	defer body.Close()

	podLogs, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(podLogs), nil
}

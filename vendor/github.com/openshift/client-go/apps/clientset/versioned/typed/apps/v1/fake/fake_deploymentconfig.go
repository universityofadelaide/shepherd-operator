// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	appsv1 "github.com/openshift/api/apps/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDeploymentConfigs implements DeploymentConfigInterface
type FakeDeploymentConfigs struct {
	Fake *FakeAppsV1
	ns   string
}

var deploymentconfigsResource = schema.GroupVersionResource{Group: "apps.openshift.io", Version: "v1", Resource: "deploymentconfigs"}

var deploymentconfigsKind = schema.GroupVersionKind{Group: "apps.openshift.io", Version: "v1", Kind: "DeploymentConfig"}

// Get takes name of the deploymentConfig, and returns the corresponding deploymentConfig object, and an error if there is any.
func (c *FakeDeploymentConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *appsv1.DeploymentConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(deploymentconfigsResource, c.ns, name), &appsv1.DeploymentConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*appsv1.DeploymentConfig), err
}

// List takes label and field selectors, and returns the list of DeploymentConfigs that match those selectors.
func (c *FakeDeploymentConfigs) List(ctx context.Context, opts v1.ListOptions) (result *appsv1.DeploymentConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(deploymentconfigsResource, deploymentconfigsKind, c.ns, opts), &appsv1.DeploymentConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &appsv1.DeploymentConfigList{ListMeta: obj.(*appsv1.DeploymentConfigList).ListMeta}
	for _, item := range obj.(*appsv1.DeploymentConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested deploymentConfigs.
func (c *FakeDeploymentConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(deploymentconfigsResource, c.ns, opts))

}

// Create takes the representation of a deploymentConfig and creates it.  Returns the server's representation of the deploymentConfig, and an error, if there is any.
func (c *FakeDeploymentConfigs) Create(ctx context.Context, deploymentConfig *appsv1.DeploymentConfig, opts v1.CreateOptions) (result *appsv1.DeploymentConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(deploymentconfigsResource, c.ns, deploymentConfig), &appsv1.DeploymentConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*appsv1.DeploymentConfig), err
}

// Update takes the representation of a deploymentConfig and updates it. Returns the server's representation of the deploymentConfig, and an error, if there is any.
func (c *FakeDeploymentConfigs) Update(ctx context.Context, deploymentConfig *appsv1.DeploymentConfig, opts v1.UpdateOptions) (result *appsv1.DeploymentConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(deploymentconfigsResource, c.ns, deploymentConfig), &appsv1.DeploymentConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*appsv1.DeploymentConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDeploymentConfigs) UpdateStatus(ctx context.Context, deploymentConfig *appsv1.DeploymentConfig, opts v1.UpdateOptions) (*appsv1.DeploymentConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(deploymentconfigsResource, "status", c.ns, deploymentConfig), &appsv1.DeploymentConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*appsv1.DeploymentConfig), err
}

// Delete takes name of the deploymentConfig and deletes it. Returns an error if one occurs.
func (c *FakeDeploymentConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(deploymentconfigsResource, c.ns, name, opts), &appsv1.DeploymentConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDeploymentConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(deploymentconfigsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &appsv1.DeploymentConfigList{})
	return err
}

// Patch applies the patch and returns the patched deploymentConfig.
func (c *FakeDeploymentConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *appsv1.DeploymentConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(deploymentconfigsResource, c.ns, name, pt, data, subresources...), &appsv1.DeploymentConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*appsv1.DeploymentConfig), err
}

// Instantiate takes the representation of a deploymentRequest and creates it.  Returns the server's representation of the deploymentConfig, and an error, if there is any.
func (c *FakeDeploymentConfigs) Instantiate(ctx context.Context, deploymentConfigName string, deploymentRequest *appsv1.DeploymentRequest, opts v1.CreateOptions) (result *appsv1.DeploymentConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateSubresourceAction(deploymentconfigsResource, deploymentConfigName, "instantiate", c.ns, deploymentRequest), &appsv1.DeploymentConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*appsv1.DeploymentConfig), err
}

// Rollback takes the representation of a deploymentConfigRollback and creates it.  Returns the server's representation of the deploymentConfig, and an error, if there is any.
func (c *FakeDeploymentConfigs) Rollback(ctx context.Context, deploymentConfigName string, deploymentConfigRollback *appsv1.DeploymentConfigRollback, opts v1.CreateOptions) (result *appsv1.DeploymentConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateSubresourceAction(deploymentconfigsResource, deploymentConfigName, "rollback", c.ns, deploymentConfigRollback), &appsv1.DeploymentConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*appsv1.DeploymentConfig), err
}

// GetScale takes name of the deploymentConfig, and returns the corresponding scale object, and an error if there is any.
func (c *FakeDeploymentConfigs) GetScale(ctx context.Context, deploymentConfigName string, options v1.GetOptions) (result *v1beta1.Scale, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetSubresourceAction(deploymentconfigsResource, c.ns, "scale", deploymentConfigName), &v1beta1.Scale{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Scale), err
}

// UpdateScale takes the representation of a scale and updates it. Returns the server's representation of the scale, and an error, if there is any.
func (c *FakeDeploymentConfigs) UpdateScale(ctx context.Context, deploymentConfigName string, scale *v1beta1.Scale, opts v1.UpdateOptions) (result *v1beta1.Scale, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(deploymentconfigsResource, "scale", c.ns, scale), &v1beta1.Scale{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Scale), err
}

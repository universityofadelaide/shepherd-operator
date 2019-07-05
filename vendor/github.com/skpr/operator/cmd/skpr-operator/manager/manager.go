package manager

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/skpr/operator/pkg/apis"
)

// New returns a configured Manager for controllers.
func New() (manager.Manager, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config")
	}

	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create manager")
	}

	err = apis.AddToScheme(mgr.GetScheme())
	if err != nil {
		return nil, errors.Wrap(err, "failed to add to scheme")
	}

	return mgr, nil
}

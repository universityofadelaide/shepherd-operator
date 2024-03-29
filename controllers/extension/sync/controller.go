package sync

import (
	osv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/universityofadelaide/shepherd-operator/controllers/extension/sync/backup"
	"github.com/universityofadelaide/shepherd-operator/controllers/extension/sync/restore"
)

// Params used by sync controllers.
type Params struct {
	Restore restore.Params
	Backup  backup.Params
}

// SetupWithManager attaches our controller to the manager.
func SetupWithManager(mgr ctrl.Manager, params Params, osclient osv1client.AppsV1Interface) error {
	if err := (&restore.Reconciler{
		Client:    mgr.GetClient(),
		OpenShift: osclient,
		Recorder:  mgr.GetEventRecorderFor(restore.ControllerName),
		Scheme:    mgr.GetScheme(),
		Params:    params.Restore,
	}).SetupWithManager(mgr); err != nil {
		return err
	}

	if err := (&backup.Reconciler{
		Client:   mgr.GetClient(),
		Recorder: mgr.GetEventRecorderFor(backup.ControllerName),
		Scheme:   mgr.GetScheme(),
		Params:   params.Backup,
	}).SetupWithManager(mgr); err != nil {
		return err
	}

	return nil
}

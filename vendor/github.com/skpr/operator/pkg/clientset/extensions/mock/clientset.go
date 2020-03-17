package mock

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	clientset "github.com/skpr/operator/pkg/clientset/extensions"
	"github.com/skpr/operator/pkg/clientset/extensions/backup"
	backupmock "github.com/skpr/operator/pkg/clientset/extensions/backup/mock"
	"github.com/skpr/operator/pkg/clientset/extensions/backupscheduled"
	backupscheduledmock "github.com/skpr/operator/pkg/clientset/extensions/backupscheduled/mock"
	"github.com/skpr/operator/pkg/clientset/extensions/exec"
	execmock "github.com/skpr/operator/pkg/clientset/extensions/exec/mock"
	"github.com/skpr/operator/pkg/clientset/extensions/restore"
	restoremock "github.com/skpr/operator/pkg/clientset/extensions/restore/mock"
)

// Clientset used for mocking.
type Clientset struct {
	backups          []*extensionsv1beta1.Backup
	backupscheduleds []*extensionsv1beta1.BackupScheduled
	restores         []*extensionsv1beta1.Restore
	execs            []*extensionsv1beta1.Exec
}

// New clientset for interacting with Entension objects.
func New(objects ...runtime.Object) (clientset.Interface, error) {
	clientset := &Clientset{}

	for _, object := range objects {
		gvk := object.GetObjectKind().GroupVersionKind()

		switch gvk.Kind {
		case backup.Kind:
			clientset.backups = append(clientset.backups, object.(*extensionsv1beta1.Backup))
		case backupscheduled.Kind:
			clientset.backupscheduleds = append(clientset.backupscheduleds, object.(*extensionsv1beta1.BackupScheduled))
		case restore.Kind:
			clientset.restores = append(clientset.restores, object.(*extensionsv1beta1.Restore))
		case exec.Kind:
			clientset.execs = append(clientset.execs, object.(*extensionsv1beta1.Exec))
		default:
			return nil, fmt.Errorf("cannot find client for: %s", gvk.Kind)
		}
	}

	return clientset, nil
}

// Backups within a namespace.
func (c *Clientset) Backups(namespace string) backup.Interface {
	filter := func(list []*extensionsv1beta1.Backup, namespace string) []*extensionsv1beta1.Backup {
		var filtered []*extensionsv1beta1.Backup

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &backupmock.Client{
		Namespace: namespace,
		Objects:   filter(c.backups, namespace),
	}
}

// BackupScheduleds within a namespace.
func (c *Clientset) BackupScheduleds(namespace string) backupscheduled.Interface {
	filter := func(list []*extensionsv1beta1.BackupScheduled, namespace string) []*extensionsv1beta1.BackupScheduled {
		var filtered []*extensionsv1beta1.BackupScheduled

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &backupscheduledmock.Client{
		Namespace: namespace,
		Objects:   filter(c.backupscheduleds, namespace),
	}
}

// Restores within a namespace.
func (c *Clientset) Restores(namespace string) restore.Interface {
	filter := func(list []*extensionsv1beta1.Restore, namespace string) []*extensionsv1beta1.Restore {
		var filtered []*extensionsv1beta1.Restore

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &restoremock.Client{
		Namespace: namespace,
		Objects:   filter(c.restores, namespace),
	}
}

// Execs within a namespace.
func (c *Clientset) Execs(namespace string) exec.Interface {
	filter := func(list []*extensionsv1beta1.Exec, namespace string) []*extensionsv1beta1.Exec {
		var filtered []*extensionsv1beta1.Exec

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &execmock.Client{
		Namespace: namespace,
		Objects:   filter(c.execs, namespace),
	}
}

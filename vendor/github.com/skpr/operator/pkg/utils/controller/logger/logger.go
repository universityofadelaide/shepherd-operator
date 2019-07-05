package logger

import (
	"github.com/prometheus/common/log"
	"github.com/rs/xid"
)

const (
	// KeyReconcile is used for tagging all reconcile activity.
	KeyReconcile = "reconcile"
	// KeyController for identifying which controller is reconciling.
	KeyController = "controller"
	// KeyNamespace for identifying which namespace the object resides.
	KeyNamespace = "namespace"
	// KeyName for identifying the object.
	KeyName = "name"
)

// New logger for server interactions.
func New(controller, namespace, name string) log.Logger {
	return log.With(KeyReconcile, xid.New().String()).With(KeyController, controller).With(KeyNamespace, namespace).With(KeyName, name)
}

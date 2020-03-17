package errors

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// IgnoreNotFound to avoid double handling.
func IgnoreNotFound(err error) error {
	if kerrors.IsNotFound(err) {
		return nil
	}

	return err
}

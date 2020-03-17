// +build unit

package mock

import (
	"testing"

	"github.com/skpr/operator/pkg/utils/slice"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestList(t *testing.T) {
	e := New()

	p := &corev1.Pod{}

	e.Event(p, "MOCK", "TEST", "Testing the Event() function")
	e.Eventf(p, "MOCK", "TEST", "Testing the %s function", "Eventf()")
	e.PastEventf(p, metav1.Now(), "MOCK", "TEST", "Testing the %s function", "PastEventf()")
	e.AnnotatedEventf(p, map[string]string{"foo": "bar"}, "MOCK", "TEST", "Testing the AnnotatedEventf() function")

	l := e.List()

	assert.True(t, slice.Contains(l, "Testing the Event() function"), "Event()")
	assert.True(t, slice.Contains(l, "Testing the Eventf() function"), "Eventf()")
	assert.True(t, slice.Contains(l, "Testing the PastEventf() function"), "PastEventf()")
	assert.True(t, slice.Contains(l, "Testing the AnnotatedEventf() function"), "AnnotatedEventf()")
}

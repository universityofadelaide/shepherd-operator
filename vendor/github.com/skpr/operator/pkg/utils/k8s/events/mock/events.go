package mock

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Mock which will store events.
type Mock struct {
	events []Event
}

// Event which occurs during an operators reconcile loop.
type Event struct {
	Timestamp   metav1.Time
	Annotations map[string]string
	Object      runtime.Object
	Type        string
	Reason      string
	Message     string
}

// Event satisfies the record.EventRecorder interface.
func (m *Mock) Event(object runtime.Object, eventtype, reason, message string) {
	m.events = append(m.events, Event{
		Timestamp: metav1.Now(),
		Object:    object,
		Type:      eventtype,
		Reason:    reason,
		Message:   message,
	})
}

// Eventf satisfies the record.EventRecorder interface.
func (m *Mock) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	m.events = append(m.events, Event{
		Timestamp: metav1.Now(),
		Object:    object,
		Type:      eventtype,
		Reason:    reason,
		Message:   fmt.Sprintf(messageFmt, args...),
	})
}

// PastEventf satisfies the record.EventRecorder interface.
func (m *Mock) PastEventf(object runtime.Object, timestamp metav1.Time, eventtype, reason, messageFmt string, args ...interface{}) {
	m.events = append(m.events, Event{
		Timestamp: timestamp,
		Object:    object,
		Type:      eventtype,
		Reason:    reason,
		Message:   fmt.Sprintf(messageFmt, args...),
	})
}

// AnnotatedEventf satisfies the record.EventRecorder interface.
func (m *Mock) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	m.events = append(m.events, Event{
		Timestamp:   metav1.Now(),
		Annotations: annotations,
		Object:      object,
		Type:        eventtype,
		Reason:      reason,
		Message:     fmt.Sprintf(messageFmt, args...),
	})
}

// List all events which have been received.
func (m *Mock) List() []string {
	var events []string

	for _, event := range m.events {
		events = append(events, event.Message)
	}

	return events
}

// New event mock.
func New() *Mock {
	return &Mock{}
}

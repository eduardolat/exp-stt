// Package eventbus provides a centralized pub/sub event bus for real-time
// communication between application components. It works as a signal system
// where publishers notify subscribers of changes, and subscribers fetch
// updated data through their own mechanisms.
package eventbus

import (
	evbus "github.com/asaskevich/EventBus"
)

// Event topics for pub/sub communication (internal use only).
const (
	topicStateChanged    = "state:changed"
	topicSettingsChanged = "settings:changed"
)

// EventBus wraps the underlying event bus and provides type-safe methods.
type EventBus struct {
	bus evbus.Bus
}

// New creates a new EventBus instance.
func New() *EventBus {
	return &EventBus{
		bus: evbus.New(),
	}
}

// PublishStateChanged notifies all subscribers that the application state has changed.
func (eb *EventBus) PublishStateChanged() {
	eb.bus.Publish(topicStateChanged)
}

// PublishSettingsChanged notifies all subscribers that settings have been updated.
func (eb *EventBus) PublishSettingsChanged() {
	eb.bus.Publish(topicSettingsChanged)
}

// SubscribeStateChanged registers a callback for state change events.
func (eb *EventBus) SubscribeStateChanged(fn func()) error {
	return eb.bus.Subscribe(topicStateChanged, fn)
}

// UnsubscribeStateChanged removes a callback from state change events.
func (eb *EventBus) UnsubscribeStateChanged(fn func()) error {
	return eb.bus.Unsubscribe(topicStateChanged, fn)
}

// SubscribeSettingsChanged registers a callback for settings change events.
func (eb *EventBus) SubscribeSettingsChanged(fn func()) error {
	return eb.bus.Subscribe(topicSettingsChanged, fn)
}

// UnsubscribeSettingsChanged removes a callback from settings change events.
func (eb *EventBus) UnsubscribeSettingsChanged(fn func()) error {
	return eb.bus.Unsubscribe(topicSettingsChanged, fn)
}

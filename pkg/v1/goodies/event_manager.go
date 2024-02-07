package goodies

import (
	"sync"
)

// EventManager is a simple event manager that allows subscribing to events
// and notifying them, so that multiple parts of the code can be notified
// when an event happens and act accordingly.
type EventManager struct {
	mu        sync.RWMutex
	listeners map[string]map[int]EventHandler
	nextID    int
}

// EventHandler is the protocol for handling events, use this to implement
// custom event handling logics for your application.
type EventHandler func(interface{})

// NewEventManager creates a new event manager.
//
// Example:
//
//	eventManager := goodies.NewEventManager()
//	eventManager.Subscribe("testEvent", func(data interface{}) {
//		fmt.Println("Event happened:", data)
//	})
//	eventManager.Notify("testEvent", "testData")
func NewEventManager() *EventManager {
	return &EventManager{
		listeners: make(map[string]map[int]EventHandler),
	}
}

// Subscribe subscribes to an event type and returns a subscription ID, use
// this ID to unsubscribe from the event later.
//
// Example:
//
//	subID := eventManager.Subscribe("testEvent", func(data interface{}) {
//		fmt.Println("Event happened:", data)
//	})
//	fmt.Println("Subscribed with ID:", subID)
func (em *EventManager) Subscribe(eventType string, handler EventHandler) int {
	em.mu.Lock()
	defer em.mu.Unlock()

	if em.listeners[eventType] == nil {
		em.listeners[eventType] = make(map[int]EventHandler)
	}

	em.nextID++
	id := em.nextID
	em.listeners[eventType][id] = handler
	return id
}

// Unsubscribe unsubscribes from an event type using its subscription ID
//
// Example:
//
//	subID := eventManager.Subscribe("testEvent", func(data interface{}) {
//		fmt.Println("Event happened:", data)
//	})
//	fmt.Println("Subscribed with ID:", subID)
//	eventManager.Unsubscribe("testEvent", subID)
func (em *EventManager) Unsubscribe(eventType string, id int) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if handlers, ok := em.listeners[eventType]; ok {
		delete(handlers, id)
	}
}

// Notify notifies an event type with some data, all the handlers subscribed
// to the event type will be called with the data.
//
// Example:
//
//	eventManager.Notify("testEvent", "testData")
func (em *EventManager) Notify(eventType string, data interface{}) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if handlers, ok := em.listeners[eventType]; ok {
		for _, handler := range handlers {
			handler(data)
		}
	}
}

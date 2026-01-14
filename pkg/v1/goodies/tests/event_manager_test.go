package tests

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/goodies"
)

func TestEventManager(t *testing.T) {
	eventManager := goodies.NewEventManager()

	// The handlerCalled variable will be used to check how many times the
	// handler is called
	handlerCalled := 0

	// Let's create a custom event type and a handler and subscribe to it
	eventType := "testEvent"
	handler := func(data interface{}) {
		// First we ensure the data is what we expect
		if data != "testData" {
			t.Errorf("Expected data to be 'testData', got %v", data)
		} else {
			t.Logf("Data is correct: %v", data)
		}
		handlerCalled++
	}
	subscriptionID := eventManager.Subscribe(eventType, handler)

	// Notify the event, expecting the handler to be called
	eventManager.Notify(eventType, "testData")
	if handlerCalled != 1 {
		t.Errorf("Expected handler to be called once, got %d", handlerCalled)
	} else {
		t.Logf("Handler was called once")
	}

	// Unsubscribe the handler to check if it gets called again on the
	// next event notification, then notify again
	eventManager.Unsubscribe(eventType, subscriptionID)
	eventManager.Notify(eventType, "testData")
	if handlerCalled > 1 {
		t.Errorf("Handler was called after unsubscribe, total calls: %d", handlerCalled)
	} else {
		t.Logf("Handler was not called after unsubscribe as expected")
	}
}

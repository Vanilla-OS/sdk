package notification

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/vanilla-os/sdk/pkg/v1/notification/types"
)

// SendNotification sends a desktop notification. It requires a Notification
// instance as an argument, you can create one using the types.NewNotification
// function.
//
// Example:
//
//	notification := types.NewNotification(
//		"BatmanApp",
//		"Batman Alert",
//		"Joker is attacking Gotham City!",
//		"batman",
//		5000,
//		types.NotificationAction{
//			Label: "Save Gotham",
//			Callback: func() {
//				fmt.Printf("Gotham saved!")
//			},
//		},
//	)
//
//	err := notification.Send()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func SendNotification(notification *types.Notification) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	var id uint32
	err = obj.Call("org.freedesktop.Notifications.Notify", 0,
		notification.AppName,                // app_name
		uint32(0),                           // replaces_id
		notification.Icon,                   // app_icon
		notification.Title,                  // summary
		notification.Message,                // body
		[]string{notification.Action.Label}, // action
		map[string]dbus.Variant{},           // hints
		notification.Timeout,                // expire_timeout
	).Store(&id)

	if err != nil {
		return err
	}

	if notification.Action.Callback == nil {
		return nil
	}

	// We need to listen for the ActionInvoked signal to know if the user
	// clicked on the notification action
	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	rule := fmt.Sprintf("member='ActionInvoked',path='/org/freedesktop/Notifications',arg0='%d'", id)
	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule)

	// Waiting for the signal
	<-c

	// Remove the match and release the callback
	conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, rule)
	notification.Action.Callback()

	return nil
}

// IsAvailable checks if the machine has notification support
// and returns true if it does, false otherwise.
//
// Example:
//
//	if notification.IsAvailable() == false {
//		fmt.Printf("Desktop notifications are not available")
//		return
//	}
func IsAvailable() bool {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false
	}

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	var capabilities []string
	err = obj.Call("org.freedesktop.Notifications.GetCapabilities", 0).Store(&capabilities)
	if err != nil {
		return false
	}

	for _, capability := range capabilities {
		// Checking if the notification server supports actions
		if capability == "actions" {
			return true
		}
	}

	return false
}

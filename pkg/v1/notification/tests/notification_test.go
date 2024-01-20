package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/notification"
	"github.com/vanilla-os/sdk/pkg/v1/notification/types"
)

func TestSendNotification(t *testing.T) {
	if notification.IsAvailable() == false {
		t.Skip("Desktop notifications are not available")
		return
	}

	notificationObj := types.NewNotification(
		"BatmanApp",
		"Batman Alert",
		"Joker is attacking Gotham City!",
		"batman",
		5000,
		types.NewNotificationAction(
			"Save Gotham",
			func() { t.Logf("Gotham saved!") },
		),
	)

	err := notification.SendNotification(notificationObj)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}

	t.Logf("Notification sent successfully")
}

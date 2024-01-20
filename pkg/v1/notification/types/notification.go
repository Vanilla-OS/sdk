package types

// Notification represents a desktop notification.
type Notification struct {
	AppName string
	Title   string
	Message string
	Icon    string
	Timeout int32
	Action  NotificationAction
}

// NotificationAction represents a desktop notification button.
type NotificationAction struct {
	Label    string
	Callback func()
}

// NewNotification creates a new Notification instance.
func NewNotification(appName, title, message, icon string, timeout int32, action NotificationAction) *Notification {
	return &Notification{
		AppName: appName,
		Title:   title,
		Message: message,
		Icon:    icon,
		Timeout: timeout,
		Action:  action,
	}
}

// NewNotificationAction creates a new NotificationAction instance.
func NewNotificationAction(label string, callback func()) NotificationAction {
	return NotificationAction{
		Label:    label,
		Callback: callback,
	}
}

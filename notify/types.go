package notify

import (
	"context"
)

type NotificationMessage interface {
	ToJSON() (map[string]interface{}, error)
}

type Notifier interface {
	Send(context.Context, NotificationMessage) (bool, error)
}

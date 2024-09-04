package notify

import (
	"context"
)

// BasicNotificationMessage 是一个泛型接口，Message 可以返回不同类型
type BasicNotificationMessage[T any] interface {
	Message() (T, error)
}

// Notifier 是一个发送通知的接口，T 是 BasicNotificationMessage 的返回类型
type Notifier[T any] interface {
	Send(ctx context.Context, message BasicNotificationMessage[T]) (bool, error)
}

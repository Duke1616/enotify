package notify

import "context"

type NotifierWrap interface {
	Send(ctx context.Context) (bool, error)
}

type NotifierWrapper[T any] struct {
	Notifier Notifier[T]
	Message  BasicNotificationMessage[T]
}

func (w NotifierWrapper[T]) Send(ctx context.Context) (bool, error) {
	return w.Notifier.Send(ctx, w.Message)
}

func WrapNotifier[T any](notifier Notifier[T], message BasicNotificationMessage[T]) NotifierWrap {
	return NotifierWrapper[T]{Notifier: notifier, Message: message}
}

type SystemNotify struct {
	Notify []NotifierWrap
}

func NewSystemNotify() *SystemNotify {
	return &SystemNotify{}
}

func (n *SystemNotify) AddNotify(notifier NotifierWrap) {
	n.Notify = append(n.Notify, notifier)
}

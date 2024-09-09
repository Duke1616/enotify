package notify

import "context"

type NotifierWrap interface {
	Send(ctx context.Context) (bool, error)
}

// GeneratedMessage 兼容动态数据
type GeneratedMessage[T any] struct {
	Data T
}

func (m *GeneratedMessage[T]) Message() (T, error) {
	return m.Data, nil
}

type NotifierWrapper[T any] struct {
	Notifier   Notifier[T]
	Message    BasicNotificationMessage[T]
	MessageGen func() (BasicNotificationMessage[T], error)
}

func (w NotifierWrapper[T]) Send(ctx context.Context) (bool, error) {
	var messageData T
	var mg BasicNotificationMessage[T]
	var err error
	if w.MessageGen != nil {
		mg, err = w.MessageGen()
		if err != nil {
			return false, err
		}

		messageData, err = mg.Message()
		if err != nil {
			return false, err
		}
	} else {
		messageData, err = w.Message.Message()
		if err != nil {
			return false, err
		}
	}

	return w.Notifier.Send(ctx, &GeneratedMessage[T]{Data: messageData})
}

// WrapNotifierStatic 方法，支持静态和动态消息生成
func WrapNotifierStatic[T any](notifier Notifier[T], message BasicNotificationMessage[T]) NotifierWrap {
	wrapper := &NotifierWrapper[T]{
		Notifier: notifier,
		Message:  message,
	}

	return wrapper
}

func WrapNotifierDynamic[T any](notifier Notifier[T], messageGen func() (
	BasicNotificationMessage[T], error)) NotifierWrap {
	wrapper := &NotifierWrapper[T]{
		Notifier:   notifier,
		MessageGen: messageGen,
	}

	return wrapper
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

//type Option[T any] func(*NotifierWrapper[T])
//
//// WithDynamicMessage 提供动态消息生成的选项函数
//func WithDynamicMessage[T any](messageGen func() (BasicNotificationMessage[T], error)) Option[T] {
//	return func(w *NotifierWrapper[T]) {
//		w.MessageGen = messageGen
//	}
//}
//
//func WithStaticMessage[T any](message BasicNotificationMessage[T]) Option[T] {
//	return func(w *NotifierWrapper[T]) {
//		w.Message = message
//	}
//}

//// WrapNotifierStatic 方法，支持静态和动态消息生成
//func WrapNotifierStatic[T any](notifier Notifier[T], message BasicNotificationMessage[T], opts ...Option[T]) NotifierWrap {
//	wrapper := &NotifierWrapper[T]{
//		Notifier: notifier,
//		Message:  message,
//	}
//
//	// 应用所有的选项函数
//	for _, opt := range opts {
//		opt(wrapper)
//	}
//	return wrapper
//}

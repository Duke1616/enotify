package message

type Content interface {
	Builder() (string, error)
	MsgType() string
}

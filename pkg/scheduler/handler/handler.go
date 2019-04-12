package handler

type MessageHandler interface {
	Name() string
	Handle(b []byte) error
}

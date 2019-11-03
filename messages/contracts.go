package messages

type BaseReply interface {
	Status() int
}

type BaseRequest interface {
	HandlerKey() string      // handlerKey
	Confirm() chan BaseReply // confirmation channel
	Respond(reply BaseReply)
}

type MessageHandler interface {
	Handle(request BaseRequest) error
}

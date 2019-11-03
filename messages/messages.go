package messages

type BaseReplyEntity struct {
	status int
}

type BaseRequestEntity struct {
	handlerKey string
	confirm    chan BaseReply
}

func (entity *BaseReplyEntity) Status() int {
	return entity.status
}

func (entity *BaseRequestEntity) HandlerKey() string {
	return entity.handlerKey
}
func (entity *BaseRequestEntity) Confirm() chan BaseReply {
	return entity.confirm
}
func (entity *BaseRequestEntity) Respond(reply BaseReply) {
	defer close(entity.confirm)
	entity.confirm <- reply
}

func MakeBaseRequestEntity(key string) *BaseRequestEntity {

	v := BaseRequestEntity{
		confirm:    make(chan BaseReply, 1), // cannot be unbuffered, because we may not yet have a listener
		handlerKey: key,
	}
	return &v
}

package messages

type FailRequest interface {
	BaseRequest
}

type FailReply interface {
	BaseReply
	Reason() string
}

type FailReplyEntity struct {
	BaseReplyEntity
	reason string
}

type FailRequestEntity struct {
	BaseRequestEntity
}

func (entity *FailReplyEntity) Reason() string {
	return entity.reason
}

func MakeFailRequestEntity() *FailRequestEntity {
	v := FailRequestEntity{
		BaseRequestEntity: *MakeBaseRequestEntity("fail-message"),
	}
	return &v
}
func MakeFailRequest() FailRequest {
	request := MakeFailRequestEntity()
	return request
}

func MakeFailReplyEntity(status int, reason string) *FailReplyEntity {
	v := FailReplyEntity{
		reason: reason,
		BaseReplyEntity: BaseReplyEntity{
			status: status,
		},
	}
	return &v
}

func MakeFailReply(status int, reason string) FailReply {
	reply := MakeFailReplyEntity(status, reason)
	return reply
}

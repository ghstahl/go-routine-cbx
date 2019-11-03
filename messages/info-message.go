package messages

type InfoRequestHandler interface {
	BaseRequest
	Name() string
}
type InfoRequestHandlers interface {
	BaseRequest
}

type InfoHandlersReply interface {
	BaseReply
	Records() []HandlerInfoRecord
}

type InfoHandlerReply interface {
	BaseReply
	Record() HandlerInfoRecord
}

type HandlerInfoRecord struct {
	Name string
	Ops  uint64
}
type InfoReplyHandlerEntity struct {
	BaseReplyEntity
	record HandlerInfoRecord
}
type InfoReplyHandlersEntity struct {
	BaseReplyEntity
	records []HandlerInfoRecord
}

type InfoRequestHandlerEntity struct {
	BaseRequestEntity
	name string
}
type InfoRequestHandlersEntity struct {
	BaseRequestEntity
}

func (entity *InfoReplyHandlerEntity) Record() HandlerInfoRecord {
	return entity.record
}
func (entity *InfoReplyHandlersEntity) Records() []HandlerInfoRecord {
	return entity.records
}

func (entity *InfoRequestHandlerEntity) Name() string {
	return entity.name
}
func MakeInfoRequestHandlerEntity(handlerKey string, name string) *InfoRequestHandlerEntity {
	v := InfoRequestHandlerEntity{
		name:              name,
		BaseRequestEntity: *MakeBaseRequestEntity(handlerKey),
	}
	return &v
}
func MakeInfoRequestHandlersEntity(handlerKey string) *InfoRequestHandlersEntity {
	v := InfoRequestHandlersEntity{
		BaseRequestEntity: *MakeBaseRequestEntity(handlerKey),
	}
	return &v
}

func MakeInfoRequestHandler(handlerKey string, name string) InfoRequestHandler {
	request := MakeInfoRequestHandlerEntity(handlerKey, name)
	return request
}

func MakeInfoRequestHandlers(handlerKey string) InfoRequestHandlers {
	request := MakeInfoRequestHandlersEntity(handlerKey)
	return request
}

func MakeInfoHandlersReplyEntity(status int, records []HandlerInfoRecord) *InfoReplyHandlersEntity {
	v := InfoReplyHandlersEntity{
		records: records,
		BaseReplyEntity: BaseReplyEntity{
			status: status,
		},
	}
	return &v
}
func MakeInfoHandlerReplyEntity(status int, record HandlerInfoRecord) *InfoReplyHandlerEntity {
	v := InfoReplyHandlerEntity{
		record: record,
		BaseReplyEntity: BaseReplyEntity{
			status: status,
		},
	}
	return &v
}

func MakeInfoHandlerReply(status int, record HandlerInfoRecord) InfoHandlerReply {
	reply := MakeInfoHandlerReplyEntity(status, record)
	return reply
}
func MakeInfoHandlersReply(status int, records []HandlerInfoRecord) InfoHandlersReply {
	reply := MakeInfoHandlersReplyEntity(status, records)
	return reply
}

package cbx

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/ghstahl/go-routine-cbx/messages"
)

type cacheContainer struct {
	Cache interface{}
}
type HanlderContainer struct {
	handler messages.MessageHandler
	ops     uint64
	cache   atomic.Value
}

const (
	_INFO_KEY = "82580e14-7d6b-4098-a0fb-a15c662fd7b2"
)

func InfoHandlerKey() string {
	return _INFO_KEY
}
func NewHandlerContainer(handler messages.MessageHandler) *HanlderContainer {
	container := &HanlderContainer{
		handler: handler,
		ops:     0,
	}
	var cc cacheContainer
	container.cache.Store(cc)

	return container
}
func (container *HanlderContainer) IncrementOps() {
	atomic.AddUint64(&container.ops, 1)
}
func (container *HanlderContainer) GetOps() uint64 {
	return atomic.LoadUint64(&container.ops)
}

func (container *HanlderContainer) UpsertCache(reply interface{}) {
	var cc cacheContainer
	cc.Cache = reply
	container.cache.Store(cc)
}
func (container *HanlderContainer) GetCache() interface{} {
	data := container.cache.Load().(cacheContainer)
	return data.Cache
}

type Runner interface {
	IssueRequest(request messages.BaseRequest, wg *sync.WaitGroup)
	AddMessageHandler(key string, handler messages.MessageHandler)
	RemoveMessageHandler(key string)
	FetchContainer(key string) (container *HanlderContainer, ok bool)
	GetMessageHandlers() *sync.Map
}

type runnerEntity struct {
	factor          int
	baseRequests    chan messages.BaseRequest
	messageHandlers sync.Map
}

func NewRunner(factor int) Runner {
	h := &runnerEntity{
		factor: factor,
	}
	h.initializeRequestChannel()
	var runner Runner
	runner = h

	rih := &runnerInfoMessageHandler{runner: runner}
	h.AddMessageHandler(InfoHandlerKey(), rih)
	return h
}

// the size of the channel queue is a factor of the numCPU.  i.e. runtime.NumCPU() * factor
func (runner *runnerEntity) initializeRequestChannel() {
	if runner.factor <= 0 || runner.factor > 100 {
		panic(fmt.Sprintf("factor out of range: 0 to %v", 100))
	}
	numCPU := runtime.NumCPU()
	queueSize := numCPU * runner.factor

	fmt.Printf("NumCPU: %v, queueSize: %v", numCPU, queueSize)

	runner.baseRequests = make(chan messages.BaseRequest, queueSize)

}

func (runner *runnerEntity) IssueRequest(request messages.BaseRequest, wg *sync.WaitGroup) {
	wg.Add(1)
	runner.baseRequests <- request
	go func() {
		defer wg.Done()
		request := <-runner.baseRequests
		runner.routeRequest(request)
	}()
}

func (runner *runnerEntity) routeRequest(request messages.BaseRequest) {
	container, ok := runner.FetchContainer(request.HandlerKey())
	if ok {
		container.IncrementOps()
		err := container.handler.Handle(request)
		if err != nil {
			log.Printf("err: %s, handler:%s", err.Error(), request.HandlerKey())
		}
		return
	}
	reason := fmt.Sprintf("err: handler does not exist, handler:%s", request.HandlerKey())
	log.Println(reason)
	reply := messages.MakeFailReply(http.StatusNotImplemented, reason)
	request.Respond(reply)
}

func (runner *runnerEntity) AddMessageHandler(key string, handler messages.MessageHandler) {
	container := NewHandlerContainer(handler)
	runner.messageHandlers.Store(key, container)
}

func (runner *runnerEntity) GetMessageHandlers() *sync.Map {
	return &runner.messageHandlers
}

func (runner *runnerEntity) RemoveMessageHandler(key string) {
	runner.messageHandlers.Delete(key)
}

func (runner *runnerEntity) FetchContainer(key string) (container *HanlderContainer, ok bool) {
	r, b := runner.messageHandlers.Load(key)

	if b {
		container = r.(*HanlderContainer)
		ok = true
		return
	}
	var v *HanlderContainer
	container = v
	ok = false
	return
}

type runnerInfoMessageHandler struct {
	runner Runner
}

func (h *runnerInfoMessageHandler) GetRunner() Runner {
	return h.runner
}

func (h *runnerInfoMessageHandler) Handle(request messages.BaseRequest) error {

	switch request.(type) {
	case messages.InfoRequestHandler:
		runner := h.GetRunner()

		container, ok := runner.FetchContainer(request.HandlerKey())

		if ok {
			record := messages.HandlerInfoRecord{
				Name: request.HandlerKey(),
				Ops:  container.GetOps(),
			}

			reply := messages.MakeInfoHandlerReply(http.StatusOK, record)
			request.Respond(reply)
			return nil
		}
		return nil
	case messages.InfoRequestHandlers:
		runner := h.GetRunner()
		var records []messages.HandlerInfoRecord
		runner.GetMessageHandlers().Range(func(k, v interface{}) bool {
			container := v.(*HanlderContainer)
			record := messages.HandlerInfoRecord{
				Name: k.(string),
				Ops:  container.GetOps(),
			}
			records = append(records, record)
			return true
		})
		reply := messages.MakeInfoHandlersReply(http.StatusOK, records)
		request.Respond(reply)
		return nil

	}
	reason := fmt.Sprintf("err: handler does not exist, handler:%s", request.HandlerKey())
	log.Println(reason)
	reply := messages.MakeFailReply(http.StatusNotImplemented, reason)
	request.Respond(reply)
	return nil
}

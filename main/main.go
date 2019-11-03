package main

import (
	"fmt"
	"runtime"
	"sync"

	cbx "github.com/ghstahl/go-routine-cbx"
	"github.com/ghstahl/go-routine-cbx/messages"
)

func main() {

	fmt.Println("starting...")

	runner := cbx.NewRunner(10)

	numCPU := runtime.NumCPU()
	queueSize := numCPU * 10

	handler := cbx.NewHandler()
	handler.ListenAndServe(9777)

	var requests []messages.BaseRequest

	var wg sync.WaitGroup

	infoRequest := messages.MakeInfoRequestHandler(cbx.InfoHandlerKey(), cbx.InfoHandlerKey())
	runner.IssueRequest(infoRequest, &wg)

	infoRequest2 := messages.MakeInfoRequestHandlers(cbx.InfoHandlerKey())
	runner.IssueRequest(infoRequest2, &wg)

	wg.Wait()

	rr := <-infoRequest.Confirm()
	switch v := rr.(type) {
	case messages.InfoHandlerReply:
		var inHR messages.InfoHandlerReply
		inHR = rr.(messages.InfoHandlerReply)

		fmt.Printf("Status %v Reason %v\n", inHR.Status(), inHR.Record())
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

	rr2 := <-infoRequest2.Confirm()
	switch v := rr2.(type) {
	case messages.InfoHandlersReply:
		var inHR messages.InfoHandlersReply
		inHR = rr2.(messages.InfoHandlersReply)

		fmt.Printf("Status %v Reason %v\n", inHR.Status(), inHR.Records())
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

	for i := 0; i < queueSize; i++ {
		request := &messages.FailRequestEntity{
			BaseRequestEntity: *messages.MakeBaseRequestEntity(fmt.Sprintf("fail-%v", i))}
		requests = append(requests, request)
		runner.IssueRequest(request, &wg)
	}
	wg.Wait()

	for _, request := range requests {
		reply := <-request.Confirm()

		switch v := reply.(type) {
		case messages.FailReply:
			var failReply messages.FailReply
			failReply = reply.(messages.FailReply)

			fmt.Printf("Status %v Reason %v\n", failReply.Status(), failReply.Reason())
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}

}

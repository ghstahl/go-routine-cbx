# go-routine-cbx

This is an excercise in using go routines to handle a request and respond back via a passed in channel.

I need to research the grpc stuff to see if I can use that engine for internal pub/reply stuff.


The request channel buffer is a factor of runtime.NumCPU.  i.e. runtime.NumCPU() * factor

This is basically what a http server would do.  

Future work is to tell a given handler type to respond with errors to excercise a server that is health but its downstream dependencies can't be reached.

In this case the server would indicate NOT READY as opposed to commiting suicide 


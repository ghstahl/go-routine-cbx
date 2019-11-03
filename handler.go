package cbx

import (
	"fmt"
	"net/http"
)

type Handler interface {
	// The Handler is an http.Handler, so it can be exposed directly and handle
	// /cbx endpoints.
	http.Handler

	// CBXEndpoint is the HTTP handler for just the /cbx endpoint
	CBXEndpoint(http.ResponseWriter, *http.Request)

	ListenAndServe(port int)
}

// basicHandler is a basic Handler implementation.
type basicHandler struct {
	http.ServeMux
}

// NewHandler creates a new basic Handler
func NewHandler() Handler {
	h := &basicHandler{}
	h.Handle("/cbx", http.HandlerFunc(h.CBXEndpoint))
	return h
}

func (s *basicHandler) CBXEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("full") != "1" {
		w.Write([]byte("{}\n"))
		return
	}
	w.Write([]byte("[{'a':'ok'}]\n"))
}

func (s *basicHandler) ListenAndServe(port int) {
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), s)
}

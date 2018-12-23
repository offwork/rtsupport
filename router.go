package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// Handler wrapper client method
type Handler func(*Client, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Router specifies rules
type Router struct {
	rules   map[string]Handler
	session *r.Session
}

// NewRouter created router instance for client
func NewRouter(session *r.Session) *Router {
	return &Router{
		rules:   make(map[string]Handler),
		session: session,
	}
}

// Handle mapping message with handler
func (r *Router) Handle(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

// FindHandler handles related message
func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}

// ServeHTTP starts http service
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	client := NewClient(socket, r.FindHandler, r.session)
	defer client.Close()
	go client.Write()
	client.Read()
}

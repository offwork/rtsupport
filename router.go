package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/websocket"
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
	rules map[string]Handler
}

// NewRouter created router instance for client
func NewRouter() *Router {
	return &Router{
		rules: make(map[string]Handler),
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
	client := NewClient(socket, r.FindHandler)
	go client.Write()
	client.Read()
}

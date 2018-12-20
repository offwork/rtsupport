package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// FindHandler shared with other source code
type FindHandler func(string) (Handler, bool)

// Message represents the outgoing message object to the receiver
type Message struct {
	Name string      `json:"name"` // name of message
	Data interface{} `json:"data"` // data of message
}

// Client shared with other source code
type Client struct {
	send        chan Message    // represents the recipient's message
	socket      *websocket.Conn // websocket connection instance
	findHandler FindHandler     // generic handler method
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}
		if handler, found := client.findHandler(message.Name); found {
			handler(client, message.Data)
		}
	}
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		fmt.Printf("%#v\n", msg)
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

// NewClient created instance client for others user
func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client {
	return &Client{
		send:        make(chan Message),
		socket:      socket,
		findHandler: findHandler,
	}
}

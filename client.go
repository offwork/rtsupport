package main

import (
	"log"

	"github.com/gorilla/websocket"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// FindHandler shared with other source code
type FindHandler func(string) (Handler, bool)

// Client shared with other source code
type Client struct {
	send         chan Message    // represents the recipient's message
	socket       *websocket.Conn // websocket connection instance
	findHandler  FindHandler     // generic handler method
	session      *r.Session
	stopChannels map[int]chan bool
	id           string
	userName     string
}

// NewStopChannel blah blah...
func (client *Client) NewStopChannel(stopkey int) chan bool {
	client.StopForKey(stopkey)
	stop := make(chan bool)
	client.stopChannels[stopkey] = stop
	return stop
}

// StopForKey blah blah...
func (client *Client) StopForKey(key int) {
	if ch, found := client.stopChannels[key]; found {
		ch <- true
		delete(client.stopChannels, key)
	}
}

// Read blah blah...
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

// Write blah blah...
func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

// Close blah blah...
func (client *Client) Close() {
	for _, ch := range client.stopChannels {
		ch <- true
	}
	close(client.send)
	r.Table("user").Get(client.id).Delete().Exec(client.session)
}

// NewClient created instance client for others user
func NewClient(socket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client {
	var user User
	user.Name = "anonymous"
	res, err := r.Table("user").Insert(user).RunWrite(session)
	if err != nil {
		log.Println(err.Error())
	}
	var id string
	if len(res.GeneratedKeys) > 0 {
		id = res.GeneratedKeys[0]
	}
	return &Client{
		send:         make(chan Message),
		socket:       socket,
		findHandler:  findHandler,
		session:      session,
		stopChannels: make(map[int]chan bool),
		id:           id,
		userName:     user.Name,
	}
}

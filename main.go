package main

import (
	"net/http"
)

// Channel represents the chat channel
type Channel struct {
	ID   string `json:"id"`   // ID of channel
	Name string `json:"name"` // name of channel
}

func main() {
	router := NewRouter()

	router.Handle("channel add", addChannel)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}

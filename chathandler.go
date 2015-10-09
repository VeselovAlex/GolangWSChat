// chathandler
package main

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

const (
	SocketBufferSize = 1024 //bytes
)

var upgrader = &ws.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: SocketBufferSize,
}

//Implement http.Handler
func (room *Room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error in Room.ServeHTTP :", err)
		return
	}

	client := &Client{
		Conn: socket,
		Msg:  make(chan []byte, MsgBufferSize),
		Room: room,
	}

	room.Join <- client
	// Leave after
	defer func() { room.Leave <- client }()
	go client.Write()
	client.Read()
}

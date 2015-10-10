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
	/*Switch to WebSocket protocol*/
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in Room.ServeHTTP :", err)
		return
	}
	
	cookie, err := r.Cookie("login")
	if err != nil {
		log.Println("Error in Room.ServeHTTP :", err)
		return
	}
	name, ok := nicknames[cookie.Value]
	
	if !ok {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		log.Println("Error in Room.ServeHTTP :", "Bad cookie")
		return
	}
	
	client := &Client{
		Conn: socket,
		Msg:  make(chan *Message, MsgBufferSize),
		Room: room,
		Name: name,
	}

	room.Join <- client
	// Leave after
	defer func() { room.Leave <- client }()
	go client.Write()
	client.Read()
}

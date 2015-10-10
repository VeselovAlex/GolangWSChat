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
	log.Println("New room connection established")
	/*Switch to WebSocket protocol*/
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error serving Room:", err)
		return
	}

	name := getName(w, r)
	if len(name) == 0 {
		log.Println("Error serving Room:", "bad user data")
		errMsg := ws.FormatCloseMessage(ws.CloseUnsupportedData,
			"Can not recognise user data")
		socket.WriteMessage(ws.CloseMessage, errMsg)
		socket.Close()
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

func getName(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("login")
	if err != nil { /*No cookie*/
		log.Println("Error reading user nick:", err)
		return ""
	}

	name, ok := nicknames[cookie.Value]
	if !ok { /*Bad cookie*/
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		log.Println("Error reading user nick:", "bad cookie")
		return ""
	}
	return name
}

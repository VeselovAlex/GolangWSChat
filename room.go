// room
package main

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

const (
	SocketBufferSize = 1024 //bytes
	MsgBufferSize    = 8    // Message pointers
	UserBufferSize   = 4    // Client pointers
)

type Room struct {
	Send  chan *Message
	Join  chan *Client
	Leave chan *Client

	Users map[*Client]bool
}

func NewRoom() *Room {
	r := new(Room)
	r.Send = make(chan *Message, MsgBufferSize)
	r.Join = make(chan *Client, UserBufferSize)
	r.Leave = make(chan *Client, UserBufferSize)
	r.Users = make(map[*Client]bool)
	return r
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Leave:
			log.Println("User", client.Name, "left the channel")
			delete(r.Users, client)
			close(client.Msg)
		case client := <-r.Join:
			log.Println("User", client.Name, "joined the channel")
			r.Users[client] = true
		case msg := <-r.Send:
			//Send to each client
			for c := range r.Users {
				select {
				case c.Msg <- msg: //Client is able to receive
					// Sent
				default:
					log.Println("User", c.Name, "left the channel while sending message")
					// Bad client, close
					delete(r.Users, c)
					close(c.Msg)
				}
			}
		}
	}
}

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
		// Probably useless, delete later
		http.SetCookie(w, &http.Cookie{Name: "login", Value: "", MaxAge: -1})
		//
		log.Println("Error reading user nick:", "bad cookie")
		return ""
	}
	return name
}

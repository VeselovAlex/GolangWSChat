// room
package main

const (
	MsgBufferSize  = 8
	UserBufferSize = 4
)

type Room struct {
	Send  chan []byte
	Join  chan *Client
	Leave chan *Client

	Users map[*Client]bool
}

func NewRoom() *Room {
	r := new(Room)
	r.Send = make(chan []byte, MsgBufferSize)
	r.Join = make(chan *Client, UserBufferSize)
	r.Leave = make(chan *Client, UserBufferSize)
	r.Users = make(map[*Client]bool)
	return r
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Leave:
			delete(r.Users, client)
			close(client.Msg)
		case client := <-r.Join:
			r.Users[client] = true
		case msg := <-r.Send:
			//Send to each client
			for c := range r.Users {
				select {
				case c.Msg <- msg: //Client is able to receive
					// Sent
				default:
					// Bad client, close
					delete(r.Users, c)
					close(c.Msg)
				}
			}
		}
	}
}

// room
package main

import (
	"log"
)

const (
	MsgBufferSize  = 8
	UserBufferSize = 4
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

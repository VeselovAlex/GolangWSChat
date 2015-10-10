// client
package main

import (
	"encoding/json"
	"log"

	ws "github.com/gorilla/websocket"
)

type Client struct {
	Conn *ws.Conn
	Msg  chan *Message
	Room *Room

	Name string
}

func (c *Client) Read() {
	for {
		/*While connection is active*/
		if _ /*Type*/, msg, err := c.Conn.ReadMessage(); err == nil {
			wrap := NewMessage(c.Name, msg)
			c.Room.Send <- wrap
		} else {
			break
		}
	}
	c.Conn.Close()
	log.Println("User", c.Name, "connection closed", "(reading)")
}

func (c *Client) Write() {
	for msg := range c.Msg {
		/*Send incoming messages to remote receiver*/
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			break
		}
		err = c.Conn.WriteMessage(ws.TextMessage, jsonMsg)
		if err != nil {
			break
		}
	}
	c.Conn.Close()
	log.Println("User", c.Name, "connection closed", "(writing)")
}

// client
package main

import (
	ws "github.com/gorilla/websocket"
	"encoding/json"
)

type Client struct {
	Conn *ws.Conn
	Msg  chan *Message
	Room *Room
	
	Name string
}

func (c *Client) Read() {
	var err error
	for err == nil {
		/*While connection is active*/
		if _ /*Type*/, msg, err := c.Conn.ReadMessage(); err == nil {
			wrap := NewMessage(c.Name, msg)
			c.Room.Send <- wrap
		}
	}
	c.Conn.Close()
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
}

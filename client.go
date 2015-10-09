// client
package main

import ws "github.com/gorilla/websocket"

type Client struct {
	Conn *ws.Conn
	Msg  chan []byte
	Room *Room
}

func (c *Client) Read() {
	var err error
	for err == nil {
		/*While connection is active*/
		if _ /*Type*/, msg, err := c.Conn.ReadMessage(); err == nil {
			c.Room.Send <- msg
		}
	}
	c.Conn.Close()
}

func (c *Client) Write() {
	for msg := range c.Msg {
		/*Send incoming messages to remote receiver*/
		err := c.Conn.WriteMessage(ws.TextMessage, msg)
		if err != nil {
			break
		}
	}
	c.Conn.Close()
}

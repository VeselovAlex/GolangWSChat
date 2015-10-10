package main

import "time"

type Message struct {
	Author    string
	Content   string
	Timestamp time.Time
}

func NewMessage(author string, content []byte) *Message {
	return &Message{author, string(content), time.Now()}
}

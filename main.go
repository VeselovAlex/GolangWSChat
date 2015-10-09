// wschat project main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)
	r := NewRoom()
	http.Handle("/ws", r)
	go r.Run()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Unable to start server on :8080 :", err)
	}
}

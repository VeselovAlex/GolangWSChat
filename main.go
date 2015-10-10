// wschat project main.go
package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	ServerHost = "127.0.0.1"
	ServerPort = "8080"
	ServerAddr = ServerHost + ":" + ServerPort
)

func serveStatic(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving", r.URL.Path)
	http.ServeFile(w, r, r.URL.Path[1:]) // Omit first '/'
}

//Shitty, but works
func main() {
	/*Init log file*/
	logFile, err := os.OpenFile("log.txt", os.O_APPEND, 0)
	if err != nil {
		logFile, err = os.Create("log.txt")
		if err != nil {
			log.Fatalln("Can not init log file")
		}
	}
	defer logFile.Close()
	logFile.Write([]byte("------------------------------------------------\n"))
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	log.Println("Simple Go chat server")
	log.Println("Server init started")

	loggedAction("Serving assets", func() {
		http.HandleFunc("/scripts/", serveStatic)
		http.HandleFunc("/styles/", serveStatic)
	})

	loggedAction("Serving chat room", func() {
		room := NewRoom()
		http.Handle("/ws", room)
		go room.Run()
	})

	loggedAction("Serving login", func() {
		login := new(Login)
		http.Handle("/login", login)
	})

	loggedAction("Serving home", func() {
		home := new(Home)
		http.Handle("/", home)
	})

	log.Println("Server initialization complete")
	log.Println("Server starts on", ServerAddr)
	err = http.ListenAndServe(ServerAddr, nil)
	if err != nil {
		log.Fatal("Unable to start server on", ServerAddr, ":", err)
	}
}

func loggedAction(msg string, action func()) {
	log.Println(msg)
	action()
	log.Println("DONE")
}

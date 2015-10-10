// wschat project main.go
package main

import (
	"log"
	"net/http"
)

func serveStatic(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving", r.URL.Path)
	http.ServeFile(w, r, r.URL.Path[1:]) // Omit first '/'
}

//Shitty, but works	
func main() {
	http.HandleFunc("/scripts/", serveStatic)
	
	http.HandleFunc("/styles/", serveStatic)
	
	r := NewRoom()
	http.Handle("/ws", r)
	go r.Run()
	
	http.HandleFunc("/login", handleLogin)

	home := &HomeHandler{}
	http.Handle("/", home)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Unable to start server on :8080 :", err)
	}
}

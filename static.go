package main

import (
	"log"
	"net/http"
)

type Static struct{}

func (s *Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving static: ", r.URL.Path)
	http.ServeFile(w, r, r.URL.Path[1:]) // Omit first '/'
}

package main

import (
	"net/http"
	"code.google.com/p/go-uuid/uuid"
)

var nicknames = make(map[string]string)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != "POST" {
		w.Header().Add("Allow", "POST")
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
	}
	
	name := r.FormValue("nickname")
	
	if len(name) == 0 || []rune(name)[0] == '@' {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} 
	
	id := uuid.New()
	nicknames[id] = name
	
	http.SetCookie(w, &http.Cookie{Name: "login", Value: id})
	http.Redirect(w, r, "/", http.StatusFound)
}
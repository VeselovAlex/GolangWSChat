package main

import (
	"log"
	"net/http"

	"code.google.com/p/go-uuid/uuid"
)

var nicknames = make(map[string]string)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving /login")
	method := r.Method
	/*Only POST allowed*/
	if method != "POST" {
		w.Header().Add("Allow", "POST")
		log.Println("Bad /login request:", "method", method)
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
	}

	name := r.FormValue("nickname")

	/*Bad User Name (it can not be empty or start with '@'*/
	if len(name) == 0 || []rune(name)[0] == '@' {
		log.Printf("Bad /login request: name %q\n", name)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	id := uuid.New()
	nicknames[id] = name

	http.SetCookie(w, &http.Cookie{Name: "login", Value: id})
	http.Redirect(w, r, "/", http.StatusFound)
}

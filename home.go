// home
package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Home struct {
	TemplateHandler
}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving home")
	h.once.Do(func() {
		home := filepath.Join(TemplateDir, "home.thtml")
		t, _ := template.ParseFiles(home)
		h.templ = t
		log.Println("Home templates parced")
	})
	method := r.Method
	if method != "GET" {
		log.Println("Error serving home:", "unsupported method", method)
		w.Header().Add("Allow", "GET")
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	nick, err := r.Cookie("login")
	if err != nil {
		log.Println("Serving login page")
		err = h.templ.ExecuteTemplate(w, "login", nil)
	} else {
		log.Println("Serving room page")
		err = h.templ.ExecuteTemplate(w, "room", nick)
	}
	if err != nil {
		log.Println("Home template error", err)
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
	}
}

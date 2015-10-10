package main;

import (
	"log"
	"path/filepath"
	"net/http"
	"html/template"
	"sync"
)

const TemplateDir = "templates"

type TemplateHandler struct {
	templ *template.Template
	once sync.Once	
}

type HomeHandler struct {
	TemplateHandler	
} 

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving home")
	h.once.Do( func() {
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
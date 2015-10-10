package main;

import (
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
	h.once.Do( func() {
		home := filepath.Join(TemplateDir, "home.thtml")
		t, _ := template.ParseFiles(home)
		h.templ = t
	})
	method := r.Method
	if method != "GET" {
		w.Header().Add("Allow", "GET")
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}
	
	nick, err := r.Cookie("login")
	if err != nil {
		err = h.templ.ExecuteTemplate(w, "login", nil)
	} else {
		 err = h.templ.ExecuteTemplate(w, "room", nick)
	}
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
	}
}
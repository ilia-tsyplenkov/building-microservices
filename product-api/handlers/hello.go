package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		h.l.Println("error reading request body", err)
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Hello, "))
	w.Write(body)
}

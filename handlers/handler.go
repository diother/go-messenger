package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/diother/go-messenger/models"
)

type MessengerHandler struct {
	message          *models.Message
	broadcastChannel chan string
}

func NewMessengerHandler(message *models.Message) *MessengerHandler {
	return &MessengerHandler{
		message:          message,
		broadcastChannel: make(chan string),
	}
}

func (h *MessengerHandler) HandleMessenger(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.New("home").ParseFiles("web/index.html"))
		if err := tmpl.ExecuteTemplate(w, "home", h.message); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		h.message.Content = r.FormValue("message")

		log.Println("Broadcasting message:", h.message.Content)
		h.broadcastChannel <- h.message.Content
	}
	return
}

func (h *MessengerHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	log.Printf("Client connected: %s", clientIP)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		message := <-h.broadcastChannel
		log.Println("Received message for client:", clientIP, "Content:", message)
		fmt.Fprintf(w, "data: %s\n\n", message)
		flusher, _ := w.(http.Flusher)
		flusher.Flush()
	}
}

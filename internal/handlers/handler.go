package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/diother/go-messenger/internal/models"
)

type BroadcasterService interface {
	Broadcast(message string)
	Subscribe(chan<- string)
	Unsubscribe(chan<- string)
}

type MessengerHandler struct {
	service BroadcasterService
	message *models.Message
}

func NewMessengerHandler(service BroadcasterService, message *models.Message) *MessengerHandler {
	return &MessengerHandler{
		service: service,
		message: message,
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

		message := strings.ReplaceAll(h.message.Content, "\n", "\\n")
		h.service.Broadcast(message)
	}
	return
}

func (h *MessengerHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientChan := make(chan string)

	h.service.Subscribe(clientChan)
	defer h.service.Unsubscribe(clientChan)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case message := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

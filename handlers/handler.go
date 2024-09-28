package handlers

import (
	"html/template"
	"net/http"

	"github.com/diother/go-messenger/models"
)

type MessengerHandler struct {
	message *models.Message
}

func NewMessengerHandler(message *models.Message) *MessengerHandler {
	return &MessengerHandler{message: message}
}

func (h *MessengerHandler) HandleMessenger(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.New("home").ParseFiles("views/index.html"))
		if err := tmpl.ExecuteTemplate(w, "home", h.message); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		h.message.Content = r.FormValue("message")

		tmpl := template.Must(template.New("home").ParseFiles("views/index.html"))
		if err := tmpl.ExecuteTemplate(w, "home", h.message); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	return
}

package main

import (
	"log"
	"net/http"

	"github.com/diother/go-messenger/handlers"
	"github.com/diother/go-messenger/models"
)

func main() {
	message := models.NewMessage("placeholder text")
	messengerHandler := handlers.NewMessengerHandler(message)
	http.HandleFunc("/", messengerHandler.HandleMessenger)

	log.Println("Server listening at port 8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

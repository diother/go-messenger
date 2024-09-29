package main

import (
	"log"
	"net/http"

	"github.com/diother/go-messenger/handlers"
	"github.com/diother/go-messenger/models"
)

func main() {
	certFile := "/etc/letsencrypt/live/hintermann.go.ro/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/hintermann.go.ro/privkey.pem"

	message := models.NewMessage("")
	messengerHandler := handlers.NewMessengerHandler(message)

	serveRouteFiles("/icons/", "./web/icons")
	serveRouteFiles("/css/", "./web/css")
	serveRouteFiles("/js/", "./web/js")

	serveSingleFile("/manifest.json", "application/json", "./web/manifest.json")
	serveSingleFile("/favicon.ico", "image/x-icon", "./web/icons/favicon.ico")

	http.HandleFunc("/", messengerHandler.HandleMessenger)
	http.HandleFunc("/events", messengerHandler.HandleEvents)

	go func() {
		log.Println("Starting HTTP to HTTPS redirect server on port 8888")
		http.ListenAndServe(":8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))
	}()

	log.Println("Starting HTTPS server on port 8889")
	err := http.ListenAndServeTLS(":8889", certFile, keyFile, nil)
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %s", err)
	}
}

func serveRouteFiles(pathPrefix, dir string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle(pathPrefix, http.StripPrefix(pathPrefix, fs))
}

func serveSingleFile(pathPrefix, contentType, dir string) {
	http.HandleFunc(pathPrefix, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, dir)
	})
}

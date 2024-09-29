package services

import (
	"log"
	"sync"
)

type BroadcasterService struct {
	clients map[chan<- string]struct{}
	mu      sync.Mutex
}

func NewBroadcasterService() *BroadcasterService {
	return &BroadcasterService{
		clients: make(map[chan<- string]struct{}),
	}
}

func (b *BroadcasterService) Broadcast(message string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for clientChan := range b.clients {
		select {
		case clientChan <- message:
		default:
			log.Printf("Client is not ready for message, dropping: %s", message)
		}
	}
}

func (b *BroadcasterService) Subscribe(clientChan chan<- string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.clients[clientChan] = struct{}{}
}

func (b *BroadcasterService) Unsubscribe(clientChan chan<- string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.clients, clientChan)
	close(clientChan)
}

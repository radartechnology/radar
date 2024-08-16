package main

import (
	"log"
)

type Hub struct {
	id         string
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	stop       chan bool
	writer     *Client
}

func newHub(session string) *Hub {
	return &Hub{
		id:         session,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		stop:       make(chan bool),
		writer:     nil,
	}
}

func (h *Hub) run() {
	log.Printf("starting hub %s", h.id)

	defer func() {
		log.Printf("stopped hub %s", h.id)
	}()

	for {
		select {
		case <-h.stop:
			h.unregisterClient(h.writer)

			for client := range h.clients {
				h.unregisterClient(client)
			}

			delete(hubs, h.id)

			return
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		log.Printf("unregistering client %s", client.id.String())
		client.close()

		delete(h.clients, client)
		close(client.send)
	}
}

func (h *Hub) broadcastMessage(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

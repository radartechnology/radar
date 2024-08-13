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
	writer     *Client
}

func newHub(session string) *Hub {
	return &Hub{
		id:         session,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
	if client.id == h.writer.id {
		log.Printf("hub %s lost its writer, closing", h.id)
		delete(hubs, h.id)

		return
	}

	if _, ok := h.clients[client]; !ok {
		return
	}

	delete(h.clients, client)
	close(client.send)
	log.Printf("unregistered client %s", client.id.String())
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

package relay

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	events <-chan *Event
}

// NewHub returns a new event hub for the provided events source
func NewHub(events <-chan *Event) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		events:     events,
	}
}

func (h *Hub) route(event *Event) {
	didSend := false
	data, err := json.Marshal(event)
	if err != nil {
		log.Errorf("Failed to encode event for %s", event.AppKey)
		return
	}

	for client := range h.clients {
		if client.appKey == event.AppKey {
			log.Debugf("Sending message for %s", event.AppKey)
			client.send <- data
			didSend = true
		}
	}

	if !didSend {
		log.Warnf("Discarded event for app '%s'", event.AppKey)
	}
}

// Run starts the event routing
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case event := <-h.events:
			h.route(event)
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

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

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// events from the AWS queue
	events <-chan *Event
}

// NewHub returns a new event hub for the provided events source
func NewHub(events <-chan *Event) *Hub {
	return &Hub{
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

	eventsReceived.Inc()

	for client := range h.clients {
		if client.appKey == event.AppKey {
			log.Debugf("Sending message for %s", event.AppKey)
			client.send <- data
			didSend = true
			eventsSent.Inc()
		}
	}

	if !didSend {
		log.Warnf("Discarded event for app '%s'", event.AppKey)
		eventsDropped.Inc()
	}
}

// Run starts the event routing
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			clientsConnected.Inc()
			clientsConnectedNow.Set(float64(len(h.clients)))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				clientsDisconnected.Inc()
				clientsConnectedNow.Set(float64(len(h.clients)))
			}
		case event := <-h.events:
			h.route(event)
		}
	}
}

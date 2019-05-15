package relay

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const prefix = "relay"

var (
	eventsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prefix,
		Name:      "events_received_total",
		Help:      "The total number of handled events",
	})

	eventsDropped = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prefix,
		Name:      "events_dropped_total",
		Help:      "The total number of dropped events",
	})

	eventsSent = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prefix,
		Name:      "events_sent_total",
		Help:      "The total number of events sent to listening clients",
	})

	clientsConnected = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prefix,
		Name:      "clients_conntected_total",
		Help:      "The total number of clients connected",
	})

	clientsDisconnected = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prefix,
		Name:      "clients_disconntected_total",
		Help:      "The total number of clients disconnected",
	})

	clientsConnectedNow = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: prefix,
		Name:      "clients_conntected",
		Help:      "The total number of clients disconnected",
	})
)

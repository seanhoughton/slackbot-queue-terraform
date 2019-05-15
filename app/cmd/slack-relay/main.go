package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seanhoughton/slackbot-queue-terraform/app/pkg/relay"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("addr", ":8080", "http service address")
var queueURL = flag.String("queue", os.Getenv("QUEUE"), "AWS queue URL")
var token = flag.String("verification", os.Getenv("VERIFICATION"), "AWS verification token")

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	flag.Parse()

	ctx := context.Background()

	events := relay.Poll(ctx, *queueURL, *token)
	hub := relay.NewHub(events)
	go hub.Run()

	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/health", health)
	r.HandleFunc("/{key}/ws", func(w http.ResponseWriter, r *http.Request) {
		relay.ServeWs(hub, w, r)
	})

	srv := &http.Server{
		Handler: r,
		Addr:    *addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

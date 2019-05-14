package main

import (
	"os"
	"context"
	"flag"
	"net/http"
	"time"
	
	"github.com/seanhoughton/slackbot-queue-terraform/app/pkg/relay"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
)


var addr = flag.String("addr", ":8080", "http service address")
var queueURL = flag.String("queue", os.Getenv("QUEUE"), "AWS queue URL")
var token = flag.String("verification", os.Getenv("VERIFICATION"), "AWS verification token")


func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	events := relay.Poll(ctx, *queueURL, *token)
	hub := relay.NewHub(events)
	go hub.Run()

	r := mux.NewRouter()

	r.HandleFunc("/", serveHome)
	r.HandleFunc("/{key}/ws", func(w http.ResponseWriter, r *http.Request) {
		relay.ServeWs(hub, w, r)
	})

	srv := &http.Server{
        Handler:      r,
        Addr:         *addr,
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
	}
	
	log.Fatal(srv.ListenAndServe())
}

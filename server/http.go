package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)

var maxHubs, _ = strconv.Atoi(os.Getenv("MAX_HUBS"))
var hubs = make(map[string]*Hub)

func runHttpServer(port string) {
	http.HandleFunc("/", wsHandler)

	addr := fmt.Sprint(":", port)
	log.Printf("listening on %s", addr)

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  pingPeriod,
		WriteTimeout: pingPeriod,
		IdleTimeout:  pingPeriod,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Header.Get("session")
	if session == "" {
		http.Error(w, "session header required", http.StatusBadRequest)
		return
	}

	parsed, err := uuid.Parse(session)
	if err != nil {
		http.Error(w, "invalid session header", http.StatusBadRequest)
		return
	}

	session = parsed.String()
	token := r.Header.Get("authorization")
	isWriter := false

	if token != "" {
		if !authenticate(r.Context(), w, token) {
			return
		}

		_, ok := hubs[session]
		if ok {
			http.Error(w, "writer already exists", http.StatusConflict)
			return
		}

		if len(hubs) >= maxHubs {
			http.Error(w, "max hubs reached", http.StatusServiceUnavailable)
			return
		}

		hub := newHub(session)
		hubs[session] = hub

		go hub.run()

		isWriter = true
	}

	hub, ok := hubs[session]
	if !ok || hub == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	log.Printf("serving session %s", session)
	serveWs(hub, w, r, isWriter)
}

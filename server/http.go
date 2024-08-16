package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var maxHubs, _ = strconv.Atoi(os.Getenv("MAX_HUBS"))
var hubs = make(map[string]*Hub)

func runHttpServer(port string) {
	http.HandleFunc("/", handleHttp)

	addr := fmt.Sprint(":", port)
	log.Printf("listening on %s", addr)

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  pingPeriod,
		WriteTimeout: pingPeriod,
		IdleTimeout:  pingPeriod,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	session := r.Header.Get("session")

	if session == "" {
		session = getSubdomain(r)
		if session == "" {
			http.Error(w, "session required", http.StatusBadRequest)
			return
		}
	}

	parsed, err := uuid.Parse(session)

	if err != nil {
		http.Error(w, "invalid session", http.StatusBadRequest)
		return
	}

	session = parsed.String()
	var hub *Hub
	token := r.Header.Get("authorization")
	isWriter := false
	var hashedToken string

	if token != "" {
		hub = checkToken(w, r, session, token)
		if hub == nil {
			return
		}
		isWriter = true
		hashedToken = hub.writer.token
	} else {
		hub, ok := hubs[session]
		if !ok || hub == nil {
			http.Error(w, "session not found", http.StatusNotFound)
			return
		}

		if len(hub.clients) >= maxHubClients {
			http.Error(w, "hub is full, rejecting connection", http.StatusServiceUnavailable)
			return
		}

		hashedToken = hashToken(token)
	}

	log.Printf("serving session %s", session)
	serveWs(hub, w, r, isWriter, hashedToken)
}

func getSubdomain(r *http.Request) string {
	host := r.Host
	parts := strings.Split(host, ".")

	if len(parts) > 2 {
		return parts[0]
	}

	return ""
}

func checkToken(w http.ResponseWriter, r *http.Request, session string, token string) *Hub {
	if len(hubs) >= maxHubs {
		http.Error(w, "max hubs reached", http.StatusServiceUnavailable)
		return nil
	}

	if _, ok := hubs[session]; ok {
		http.Error(w, "writer already exists", http.StatusConflict)
		return nil
	}

	hashedToken := hashToken(token)

	if !authenticate(r.Context(), w, hashedToken) {
		return nil
	}

	hub := newHub(session)
	hubs[session] = hub

	go hub.run()

	return hub
}

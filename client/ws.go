package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

func getWsConn() *websocket.Conn {
	token, err := getCachedToken()
	if err != nil {
		token = stringPrompt("enter token: ")
		if err = cacheToken(token); err != nil {
			log.Printf("failed to cache token: %v", err)
		}
	}

	u := url.URL{Scheme: "wss", Host: *serverAddress, Path: "/"}

	header := http.Header{}
	header.Add("authorization", token)

	id, err := uuid.NewV7()
	if err != nil {
		log.Fatalf("failed to generate uuid: %v", err)
	}

	session := id.String()
	header.Add("session", session)

	log.Printf("dialing websocket: %s", u.String())
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, res, err := dialer.Dial(u.String(), header)

	if err != nil {
		buf := make([]byte, 4096)
		n, _ := res.Body.Read(buf)
		log.Fatalf("websocket dial failed: %s", buf[:n])
	}

	log.Printf("connected to websocket: %s", u.String())
	log.Printf("your radar link is http://radar.kidua.net?ip=%s.radar.technology", session)

	return conn
}

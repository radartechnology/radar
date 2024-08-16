package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
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

	u := url.URL{Scheme: "ws", Host: *serverAddress, Path: "/"}

	header := http.Header{}
	header.Add("authorization", token)

	id, err := uuid.NewV7()
	if err != nil {
		log.Printf("failed to generate uuid: %v", err)
		fatal()
	}

	session := id.String()
	header.Add("session", session)

	log.Printf("dialing websocket: %s", u.String())
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, res, err := dialer.Dial(u.String(), header)

	if err != nil {
		log.Printf("failed to dial websocket: %v", err)
		buf := make([]byte, 4096)
		n, _ := res.Body.Read(buf)
		log.Printf("websocket dial failed: %s", buf[:n])
		fatal()
	}

	log.Printf("connected to websocket: %s", u.String())
	link := fmt.Sprintf("http://radar.kidua.net?ip=%s.radar.technology", session)
	log.Printf("your radar link is %s", link)
	_ = browser.OpenURL(link)

	return conn
}

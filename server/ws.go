package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var maxHubClients, _ = strconv.Atoi(os.Getenv("MAX_HUB_CLIENTS"))

const (
	pingPeriod = 250 * time.Millisecond
	bufSize    = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		//goland:noinspection HttpUrlsUsage
		// return r.Header.Get("Origin") == "http://radar.kidua.net"
		return true
	},
}

type Client struct {
	id   uuid.UUID
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func newClient(hub *Hub, conn *websocket.Conn) (*Client, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Client{
		id:   id,
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}, nil
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request, isWriter bool) {
	log.SetPrefix(fmt.Sprintf("[%s] ", hub.id))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to upgrade connection: %v", err), http.StatusInternalServerError)
		return
	}

	client, err := newClient(hub, conn)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create client: %v", err), http.StatusInternalServerError)
		return
	}

	if isWriter {
		log.Printf("registering writer %s", client.id.String())
		hub.writer = client

		go client.read()
	} else {
		log.Printf("registering client %s", client.id.String())

		client.hub.register <- client
		go client.write()
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.close()
	}()

	log.Printf("started write worker")

	for {
		select {
		case message, ok := <-c.send:
			if err := sendMessage(c, message, ok); err != nil {
				log.Printf("failed to send message: %v", err)
				return
			}
		case <-ticker.C:
			if err := sendPing(c); err != nil {
				log.Printf("failed to send ping: %v", err)
				return
			}
		}
	}
}

func (c *Client) close() {
	c.hub.unregister <- c

	err := c.conn.Close()
	if err != nil {
		log.Printf("failed to close connection: %v", err)
		return
	}

	log.Printf("successfully closed connection and stopped worker")
}

func sendMessage(c *Client, message []byte, ok bool) error {
	_ = c.conn.SetWriteDeadline(time.Now().Add(pingPeriod))

	if !ok {
		err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		if err != nil {
			return err
		}

		return error(nil)
	}

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	defer func(w io.WriteCloser) {
		err = w.Close()
		if err != nil {
			log.Printf("failed to close writer: %v", err)
		}
	}(w)

	message = bytes.TrimSpace(message)
	_, _ = w.Write(message)

	n := len(c.send)
	for range n {
		_, _ = w.Write(<-c.send)
	}

	return nil
}

func sendPing(c *Client) error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(pingPeriod)); err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *Client) read() {
	defer c.close()
	c.conn.SetReadLimit(bufSize)

	log.Println("starting read worker")

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message: %v", err)
			break
		}

		message = bytes.TrimSpace(message)
		c.hub.broadcast <- message
	}
}

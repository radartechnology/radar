package main

import (
	"context"
	"github.com/Microsoft/go-winio"
	"github.com/gorilla/websocket"
	"log"
	"net"
)

// todo: check if this changes
const (
	pipeName = `\\.\pipe\23d339ddef616cb0a5b9d0be60a289bc4ae87cc62cfd12b8f322e6310c1eea66`
	bufSize  = 4096
)

func startTransmission() {
	pipe := getPipe()
	conn := getWsConn()

	write(conn, pipe)
}

func getPipe() net.Conn {
	log.Println("looking for data pipe")

	pipe, err := dialPipe()

	for err != nil {
		log.Fatalf("dial failed: %v", err)
	}

	log.Println("connected to data pipe")

	return pipe
}

func dialPipe() (net.Conn, error) {
	return winio.DialPipeAccess(context.Background(), pipeName, 1)
}

func write(conn *websocket.Conn, pipe net.Conn) {
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("close failed: %v", err)
		}
	}(conn)

	for {
		buf := make([]byte, bufSize)

		n, err := pipe.Read(buf)
		if err != nil {
			log.Printf("read from pipe failed: %v", err)
			return
		}

		err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			log.Printf("write to websocket failed: %v", err)
			break
		}
	}
}

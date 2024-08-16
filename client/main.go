package main

import (
	"flag"
	"os"
	"time"
)

var serverAddress = flag.String("server-address", "radar.technology:1887", "example: radar.technology:1887")

func main() {
	flag.Parse()
	startTransmission()
}

func fatal() {
	time.Sleep(5 * time.Second)
	os.Exit(1)
}

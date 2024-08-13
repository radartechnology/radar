package main

import "flag"

var serverAddress = flag.String("server-address", "radar.technology", "example: radar.technology")

func main() {
	flag.Parse()
	startTransmission()
}

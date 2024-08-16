package main

import "flag"

var serverAddress = flag.String("server-address", "radar.technology:1887", "example: radar.technology:1887")

func main() {
	flag.Parse()
	startTransmission()
}

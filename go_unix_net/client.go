package main

import (
	"log"
	"net"
)

func main() {
	client, err := net.Dial("tcp", ":3030")
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer client.Close()

	n, err := client.Write([]byte("Hello, world!"))
	if err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Printf("Sent: %d bytes", n)
}

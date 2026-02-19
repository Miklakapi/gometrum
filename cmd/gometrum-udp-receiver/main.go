package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	addr := "0.0.0.0:8222"

	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}
	defer conn.Close()

	log.Printf("UDP receiver listening on %s\n", addr)

	buf := make([]byte, 65535)

	for {
		n, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("read error: %v\n", err)
			continue
		}

		fmt.Printf("[%s] %s\n", remoteAddr.String(), string(buf[:n]))
	}
}

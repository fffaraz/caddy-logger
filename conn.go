package main

import (
	"log"
	"net"
	"os"
)

func getUnixConn(socketAddr string) (net.Conn, error) {
	if _, err := os.Stat(socketAddr); err == nil {
		log.Println("Removing existing socket")
		if err := os.Remove(socketAddr); err != nil {
			log.Println("Error removing existing socket:", err)
			return nil, err
		}
	}

	addr := net.UnixAddr{
		Name: socketAddr,
		Net:  "unixgram",
	}
	conn, err := net.ListenUnixgram("unixgram", &addr)
	if err != nil {
		log.Println("Error listening:", err)
		return nil, err
	}

	log.Println("Listening on socket", socketAddr)
	return conn, nil
}

func getUDPConn(listenPort int) (net.Conn, error) {
	addr := net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: listenPort,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Println("Error listening:", err)
		return nil, err
	}

	log.Println("Listening on port", listenPort)
	return conn, nil
}

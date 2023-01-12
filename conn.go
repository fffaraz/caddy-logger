package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"

	"gorm.io/gorm"
)

func getUnixConn(socketAddr string) (net.Conn, error) {
	if _, err := os.Stat(socketAddr); err == nil {
		fmt.Println("Removing existing socket")
		if err := os.Remove(socketAddr); err != nil {
			fmt.Println("Error removing existing socket:", err)
			return nil, err
		}
	}

	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{
		Name: socketAddr,
		Net:  "unixgram",
	})
	if err != nil {
		fmt.Println("Error listening:", err)
		return nil, err
	}
	fmt.Println("Listening on", socketAddr)

	return conn, nil
}

func getUdpConn(listenPort int) (net.Conn, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: listenPort,
	})
	if err != nil {
		fmt.Println("Error listening:", err)
		return nil, err
	}
	fmt.Println("Listening on port", listenPort)

	return conn, nil
}

func readConn(conn net.Conn, db *gorm.DB, wg *sync.WaitGroup) {
	for {
		buf := make([]byte, 1024*64)
		nr, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		_, info, err := getLogMessage(buf[:nr])
		if err != nil {
			fmt.Println("Error getting log message:", err)
			continue
		}
		if err := db.Create(info).Error; err != nil {
			fmt.Println("Error saving log message:", err)
			continue
		}
		infoStr, _ := json.Marshal(info)
		fmt.Println(string(infoStr) + "\n")
	}
	wg.Done()
}

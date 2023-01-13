package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

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

	fmt.Println("Listening on socket", socketAddr)
	return conn, nil
}

func getUdpConn(listenPort int) (net.Conn, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
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
		startTime := time.Now()
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		_, log, err := getMessage(buf[:nr])
		if err != nil {
			fmt.Println("Error getting log message:", err)
			continue
		}
		if err := db.Create(log).Error; err != nil {
			fmt.Println("Error saving log message:", err)
			continue
		}
		if false {
			infoStr, _ := json.Marshal(log)
			fmt.Println(time.Since(startTime).Milliseconds(), string(infoStr))
		}
	}
	wg.Done()
}

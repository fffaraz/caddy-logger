package main

import (
	"fmt"
	"net"
	"os"

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

	addr := net.UnixAddr{
		Name: socketAddr,
		Net:  "unixgram",
	}
	conn, err := net.ListenUnixgram("unixgram", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return nil, err
	}

	fmt.Println("Listening on socket", socketAddr)
	return conn, nil
}

func getUdpConn(listenPort int) (net.Conn, error) {
	addr := net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: listenPort,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return nil, err
	}

	fmt.Println("Listening on port", listenPort)
	return conn, nil
}

func readConn(conn net.Conn, db *gorm.DB) {
	buf := make([]byte, 1024*64)
	msg := &Message{}
	log := &Log{}
	for {
		nr, err := conn.Read(buf)
		// startTime := time.Now()
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		if err := getMessage(buf[:nr], msg, log); err != nil {
			fmt.Println("Error getting log message:", err)
			continue
		}
		if err := db.Create(log).Error; err != nil {
			fmt.Println("Error saving log message:", err)
			break
		}
		// infoStr, _ := json.Marshal(log)
		// fmt.Println(time.Since(startTime).Milliseconds(), string(infoStr))
	}
}

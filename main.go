package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	apiPort := flag.Int("a", 0, "api port")
	dbPath := flag.String("d", "", "database path")
	listenPort := flag.Int("p", 0, "listen port")
	socketAddr := flag.String("s", "", "socket address")
	flag.Parse()

	db, err := getDB(*dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	if *apiPort != 0 {
		go startApi(*apiPort, db)
	}

	var conn net.Conn
	if *socketAddr != "" {
		conn, err = getUnixConn(*socketAddr)
		if err != nil {
			fmt.Println("Error opening socket:", err)
			return
		}
	}
	if *listenPort != 0 {
		conn, err = getUdpConn(*listenPort)
		if err != nil {
			fmt.Println("Error opening port:", err)
			return
		}
	}
	if conn == nil {
		fmt.Println("No connection specified")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go readConn(conn, db, &wg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	conn.Close()
	wg.Wait()

	if *socketAddr != "" {
		os.Remove(*socketAddr)
	}

	fmt.Println("\nExiting")
}

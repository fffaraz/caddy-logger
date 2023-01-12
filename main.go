package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	dbPath := flag.String("d", "", "database path")
	listenPort := flag.Int("p", 0, "listen port")
	socketAddr := flag.String("s", "", "socket address")
	flag.Parse()

	db, err := getDB(*dbPath)
	if err != nil {
		log.Fatal(err)
	}

	var conn net.Conn
	if *socketAddr != "" {
		conn, err = getUnixConn(*socketAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *listenPort != 0 {
		conn, err = getUdpConn(*listenPort)
		if err != nil {
			log.Fatal(err)
		}
	}
	if conn == nil {
		log.Fatal("No connection")
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

	fmt.Println()
	log.Println("Exiting")
}

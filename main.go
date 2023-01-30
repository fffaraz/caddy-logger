package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	apiPort := flag.Int("a", 0, "api port")
	dbPath := flag.String("d", "db.sqlite", "database path")
	listenPort := flag.Int("p", 0, "listen port")
	socketAddr := flag.String("s", "", "socket address")
	flag.Parse()

	currentTime := time.Now().Format("2006-01-02-15-04-05")
	newDBPath := strings.Replace(*dbPath, "timestamp", currentTime, -1)

	db, err := getDB(newDBPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Error getting database:", err)
		return
	}
	defer sqlDB.Close()

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
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		readConn(conn, db)
		ch <- syscall.SIGTERM
		wg.Done()
	}()

	<-ch
	conn.Close()
	wg.Wait()

	if *socketAddr != "" {
		os.Remove(*socketAddr)
	}

	fmt.Println("\nExiting")
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	apiPort := flag.Int("api", 0, "api port")
	dbPath := flag.String("db", "db.sqlite", "database path")
	listenPort := flag.Int("port", 0, "listen port")
	socketAddr := flag.String("socket", "", "socket address")
	flag.Parse()

	currentTime := time.Now().Format("2006-01-02-15-04-05")
	newDBPath := strings.ReplaceAll(*dbPath, "timestamp", currentTime)

	db, err := getDB(newDBPath)
	if err != nil {
		log.Println("Error opening database:", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Error getting database:", err)
		return
	}
	defer sqlDB.Close()

	if *apiPort != 0 {
		go startAPI(*apiPort, db)
	}

	var conn net.Conn
	if *socketAddr != "" {
		conn, err = getUnixConn(*socketAddr)
		if err != nil {
			log.Println("Error opening socket:", err)
			return
		}
	}
	if *listenPort != 0 {
		conn, err = getUDPConn(*listenPort)
		if err != nil {
			log.Println("Error opening port:", err)
			return
		}
	}
	if conn == nil {
		log.Println("No connection specified")
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

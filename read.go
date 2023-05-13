package main

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"gorm.io/gorm"
)

var msgPerSecond float64

func readConn(conn net.Conn, db *gorm.DB) {
	buf := make([]byte, 1024*64)
	var msg Message
	var obj Log

	numMessages := 0
	totalCounter := 0
	startTime := time.Now()

	for {
		numMessages++
		totalCounter++
		if totalCounter > 5_000_000 {
			log.Println("Restarting")
			break
		}

		elapsedTime := time.Since(startTime)
		if elapsedTime >= time.Second {
			perSecond := float64(numMessages) / elapsedTime.Seconds()
			msgPerSecond = approxRollingAverage(msgPerSecond, perSecond, 30)
			numMessages = 0
			startTime = time.Now()
		}

		nr, err := conn.Read(buf)
		processTimeBegin := time.Now()
		if err != nil {
			log.Println("Error reading:", err)
			break
		}

		if err := getMessage(buf[:nr], &msg, &obj); err != nil {
			log.Println("Error getting log message:", err)
			continue
		}

		if err := db.Create(&obj).Error; err != nil {
			log.Println("Error saving log message:", err)
			break
		}

		processTimeEnd := time.Since(processTimeBegin)
		if processTimeEnd > time.Millisecond*100 {
			infoStr, err := json.Marshal(&obj)
			if err != nil {
				log.Println("Error marshalling log:", err)
				continue
			}
			log.Println("Slow query:", processTimeEnd.Milliseconds(), string(infoStr))
		}
	}
}

func approxRollingAverage(avg, input float64, n int) float64 {
	// accumulator = (alpha * new_value) + (1.0 - alpha) * accumulator
	avg -= avg / float64(n)
	avg += input / float64(n)
	return avg
}

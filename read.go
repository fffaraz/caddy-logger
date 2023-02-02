package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"gorm.io/gorm"
)

var msgPerSecond float64

func readConn(conn net.Conn, db *gorm.DB) {
	buf := make([]byte, 1024*64)
	msg := &Message{}
	log := &Log{}

	numMessages := 0
	totalCounter := 0
	startTime := time.Now()

	for {
		numMessages++
		totalCounter++
		if totalCounter > 5_000_000 {
			fmt.Println("Restarting")
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

		processTimeEnd := time.Since(processTimeBegin)
		if processTimeEnd > time.Millisecond*100 {
			infoStr, _ := json.Marshal(log)
			fmt.Println("Slow query:", processTimeEnd.Milliseconds(), string(infoStr))
		}
	}
}

func approxRollingAverage(avg, input float64, N int) float64 {
	// accumulator = (alpha * new_value) + (1.0 - alpha) * accumulator
	avg -= avg / float64(N)
	avg += input / float64(N)
	return avg
}

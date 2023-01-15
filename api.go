package main

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func startApi(port int, db *gorm.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var count int64
		db.Model(&Log{}).Count(&count)
		fmt.Fprintf(w, "Hello, %s\n%d\n", r.RemoteAddr, count)
	})

	http.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		orderBy := "host"
		switch r.URL.Query().Get("sort") {
		case "hits":
			orderBy = "hits DESC"
		case "traffic":
			orderBy = "traffic DESC"
		}
		var results []ResultHost
		tx := db.Raw("SELECT host, COUNT(id) AS hits, SUM(size) AS traffic FROM logs GROUP BY host ORDER BY " + orderBy).Scan(&results)
		if tx.Error != nil {
			fmt.Fprintf(w, "Error: %s", tx.Error)
			return
		}
		fmt.Fprintf(w, "Host\tHits\tTraffic\n")
		for _, result := range results {
			fmt.Fprintf(w, "%s\t%d\t%d\n", result.Host, result.Hits, result.Traffic)
		}
	})

	fmt.Printf("Starting API on port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil); err != nil {
		fmt.Println("API error:", err)
	}
}

type ResultHost struct {
	Host    string
	Hits    int
	Traffic int64
}

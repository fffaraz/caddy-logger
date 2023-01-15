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
		var results []resultHosts
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

	http.HandleFunc("/host", func(w http.ResponseWriter, r *http.Request) {
		host := r.URL.Query().Get("h")
		if host == "" {
			fmt.Fprintf(w, "Error: missing host (?h=...) parameter")
			return
		}
		var results []resultHost
		tx := db.Raw("SELECT date(time_stamp) AS date, COUNT(id) AS hits, SUM(size) AS traffic FROM logs WHERE host like ? OR host like ? GROUP BY date(time_stamp) ORDER BY date(time_stamp) DESC", host, "%."+host).Scan(&results)
		if tx.Error != nil {
			fmt.Fprintf(w, "Error: %s", tx.Error)
			return
		}
		fmt.Fprintf(w, "Date\tHits\tTraffic\n")
		for _, result := range results {
			fmt.Fprintf(w, "%s\t%d\t%d\n", result.Date, result.Hits, result.Traffic)
		}
	})

	fmt.Printf("Starting API on port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil); err != nil {
		fmt.Println("API error:", err)
	}
}

type resultHost struct {
	Date    string
	Hits    int
	Traffic int64
}

type resultHosts struct {
	Host    string
	Hits    int
	Traffic int64
}

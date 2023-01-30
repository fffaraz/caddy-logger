package main

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func startApi(port int, db *gorm.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var count int64
		var model Log
		db.Model(&model).Count(&count)
		fmt.Fprintf(w, "Hello, %s\n%d\n", r.RemoteAddr, count)
	})

	http.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		orderBy := "domain"
		switch r.URL.Query().Get("sort") {
		case "hits":
			orderBy = "hits DESC"
		case "traffic":
			orderBy = "traffic DESC"
		}
		var results []struct {
			Domain  string
			Hits    int
			Traffic int64
		}
		tx := db.Raw("SELECT domain, COUNT(id) AS hits, SUM(size) AS traffic FROM logs GROUP BY domain ORDER BY " + orderBy).Scan(&results)
		if tx.Error != nil {
			fmt.Fprintf(w, "Error: %s", tx.Error)
			return
		}
		fmt.Fprintf(w, "Domain\tHits\tTraffic\n")
		for _, result := range results {
			fmt.Fprintf(w, "%s\t%d\t%d\n", result.Domain, result.Hits, result.Traffic)
		}
	})

	http.HandleFunc("/domain", func(w http.ResponseWriter, r *http.Request) {
		domain := r.URL.Query().Get("d")
		if domain == "" {
			fmt.Fprintf(w, "Error: missing domain (?d=...) parameter")
			return
		}
		var results []struct {
			Date    string
			Hits    int
			Traffic int64
		}
		tx := db.Raw("SELECT date(time_stamp) AS date, COUNT(id) AS hits, SUM(size) AS traffic FROM logs WHERE domain like ? GROUP BY date(time_stamp) ORDER BY date(time_stamp) DESC", domain).Scan(&results)
		if tx.Error != nil {
			fmt.Fprintf(w, "Error: %s", tx.Error)
			return
		}
		fmt.Fprintf(w, "Date\tHits\tTraffic\n")
		for _, result := range results {
			fmt.Fprintf(w, "%s\t%d\t%d\n", result.Date, result.Hits, result.Traffic)
		}
	})

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		simple := r.URL.Query().Get("simple") == "1"
		domain := r.URL.Query().Get("d")
		if domain == "" {
			fmt.Fprintf(w, "Error: missing domain (?d=...) parameter")
			return
		}
		var results []Log
		tx := db.Where("domain like ?", domain).Order("id DESC").Limit(5000).Find(&results)
		if tx.Error != nil {
			fmt.Fprintf(w, "Error: %s", tx.Error)
			return
		}
		printLogs(w, results, simple)
	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		simple := r.URL.Query().Get("simple") == "1"
		var results []Log
		tx := db.Order("id DESC").Limit(5000).Find(&results)
		if tx.Error != nil {
			fmt.Fprintf(w, "Error: %s", tx.Error)
			return
		}
		printLogs(w, results, simple)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	})

	fmt.Printf("Starting API on port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil); err != nil {
		fmt.Println("API error:", err)
	}
}

func printLogs(w http.ResponseWriter, results []Log, simple bool) {
	if simple {
		fmt.Fprintf(w, "ID\tStatus\tRemoteIp\tDomain\tHost\tUri\tUserAgent\n")
		for _, result := range results {
			fmt.Fprintf(w, "%d\t%d\t%s\t%s\t%s\t%s\t%s\n", result.ID, result.Status, result.RemoteIp, result.Domain, result.Host, result.Uri, result.UserAgent)
		}
	} else {
		fmt.Fprintf(w, "ID\tTimeStamp\tDuration\tSize\tStatus\tRemoteIp\tRemotePort\tProto\tMethod\tHost\tDomain\tUri\tUserAgent\tCfRay\tCfConnectingIp\tCfIPCountry\tXForwardedFor\tTlsServerName\n")
		for _, result := range results {
			fmt.Fprintf(w, "%d\t%s\t%d\t%d\t%d\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", result.ID,
				result.TimeStamp, result.Duration, result.Size, result.Status, result.RemoteIp, result.RemotePort, result.Proto, result.Method,
				result.Host, result.Domain, result.Uri, result.UserAgent, result.CfRay, result.CfConnectingIp, result.CfIPCountry, result.XForwardedFor,
				result.TlsServerName)
		}
	}
}

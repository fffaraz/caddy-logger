package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// "github.com/glebarez/sqlite"

type Log struct {
	ID             uint          `gorm:"primarykey"`
	TimeStamp      time.Time     ``
	Duration       time.Duration ``
	Size           int64         ``
	Status         int           `gorm:"index"`
	RemoteIP       string        `gorm:"index"`
	RemotePort     int           ``
	Proto          string        ``
	Method         string        ``
	Host           string        ``
	Domain         string        `gorm:"index"`
	URI            string        ``
	UserAgent      string        ``
	CfRay          string        ``
	CfConnectingIP string        ``
	CfIPCountry    string        ``
	XForwardedFor  string        ``
	TLSServerName  string        ``
}

func getDB(dbPath string) (*gorm.DB, error) {
	// _pragma=journal_mode(wal)
	db, err := gorm.Open(sqlite.Open(dbPath+"?_journal_mode=wal"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Log{}); err != nil {
		return nil, err
	}

	return db, nil
}

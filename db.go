package main

import (
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// "gorm.io/driver/sqlite"

type Log struct {
	ID             uint          `gorm:"primarykey"`
	TimeStamp      time.Time     ``
	Duration       time.Duration ``
	Size           int64         ``
	Status         int           `gorm:"index"`
	RemoteIp       string        `gorm:"index"`
	RemotePort     int           ``
	Proto          string        ``
	Method         string        ``
	Host           string        `gorm:"index"`
	Uri            string        ``
	UserAgent      string        ``
	CfRay          string        `` // Cloudflare Ray ID
	CfConnectingIp string        ``
	CfIPCountry    string        ``
	XForwardedFor  string        ``
	TlsServerName  string        ``
}

func getDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Log{}); err != nil {
		return nil, err
	}

	return db, nil
}

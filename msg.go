package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func getMessage(buf []byte, msg *Message, log *Log) error {
	// fmt.Println(string(buf))

	msg.Level = ""
	msg.Ts = 0
	msg.Logger = ""
	msg.Msg = ""
	msg.UserId = ""
	msg.Duration = 0
	msg.Size = 0
	msg.Status = 0
	msg.Request.RemoteIp = ""
	msg.Request.RemotePort = ""
	msg.Request.Proto = ""
	msg.Request.Method = ""
	msg.Request.Host = ""
	msg.Request.Uri = ""
	msg.Request.Headers = nil
	msg.Request.Tls = nil
	msg.RespHeaders = nil

	if err := json.Unmarshal(buf, msg); err != nil {
		return err
	}

	log.ID = 0
	log.TimeStamp = time.Unix(int64(msg.Ts), 0)
	log.Duration = time.Duration(msg.Duration * float64(time.Second))
	log.Size = msg.Size
	log.Status = msg.Status
	log.RemoteIp = msg.Request.RemoteIp
	log.RemotePort, _ = strconv.Atoi(msg.Request.RemotePort)
	log.Proto = msg.Request.Proto
	log.Method = msg.Request.Method
	log.Host = msg.Request.Host
	log.Domain = getDomain(msg.Request.Host)
	log.Uri = msg.Request.Uri
	log.UserAgent = getHeader(msg, "User-Agent")
	log.CfRay = getHeader(msg, "CF-Ray")
	log.CfConnectingIp = getHeader(msg, "CF-Connecting-IP")
	log.CfIPCountry = getHeader(msg, "CF-IPCountry")
	log.XForwardedFor = getHeader(msg, "X-Forwarded-For")
	if msg.Request.Tls != nil {
		log.TlsServerName = msg.Request.Tls.ServerName
	} else {
		log.TlsServerName = ""
	}

	return nil
}

func getHeader(msg *Message, key string) string {
	if msg.Request.Headers[key] != nil {
		return msg.Request.Headers[key][0]
	}
	return ""
}

func getDomain(host string) string {
	host = strings.ToLower(host)

	if len(host) < 4 {
		return host
	}

	if i := strings.Index(host, "]:"); i >= 0 {
		return host[1:i] // IPv6 + port
	}

	if host[0] == '[' {
		return host[1 : len(host)-1] // IPv6
	}

	if i := strings.IndexByte(host, ':'); i >= 0 {
		host = host[:i] // remove port
	}

	parts := strings.Split(host, ".")
	if len(parts) <= 2 {
		return host
	}

	tld := parts[len(parts)-1]
	if _, err := strconv.Atoi(tld); err == nil {
		return host // IPv4
	}

	return strings.Join(parts[len(parts)-2:], ".")
}

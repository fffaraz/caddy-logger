package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func getLogMessage(buf []byte) (*LogMessage, *LogInfo, error) {
	fmt.Println(string(buf))

	msg := &LogMessage{}
	if err := json.Unmarshal(buf, msg); err != nil {
		return nil, nil, err
	}

	info := &LogInfo{}
	info.TimeStamp = time.Unix(int64(msg.Ts), 0)
	info.Duration = time.Duration(msg.Duration * float64(time.Second))
	info.Size = msg.Size
	info.Status = msg.Status
	info.RemoteIp = msg.Request.RemoteIp
	info.RemotePort, _ = strconv.Atoi(msg.Request.RemotePort)
	info.Proto = msg.Request.Proto
	info.Method = msg.Request.Method
	info.Host = msg.Request.Host
	info.Uri = msg.Request.Uri
	info.UserAgent = getHeader(msg, "User-Agent")
	info.CfConnectingIp = getHeader(msg, "Cf-Connecting-Ip")
	info.CfIpcountry = getHeader(msg, "Cf-Ipcountry")
	info.XForwardedFor = getHeader(msg, "X-Forwarded-For")
	if msg.Request.Tls != nil {
		info.TlsServerName = msg.Request.Tls.ServerName
	}

	return msg, info, nil
}

func getHeader(msg *LogMessage, key string) string {
	if msg.Request.Headers[key] != nil {
		return msg.Request.Headers[key][0]
	}
	return ""
}

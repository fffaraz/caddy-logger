package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type LogTls struct {
	Resumed     bool   `json:"resumed"`
	Version     int    `json:"version"`
	CipherSuite int    `json:"cipher_suite"`
	Proto       string `json:"proto"`
	ServerName  string `json:"server_name"`
}

type LogRequest struct {
	RemoteIp   string              `json:"remote_ip"`
	RemotePort string              `json:"remote_port"`
	Proto      string              `json:"proto"`
	Method     string              `json:"method"`
	Host       string              `json:"host"`
	Uri        string              `json:"uri"`
	Headers    map[string][]string `json:"headers"`
	Tls        *LogTls             `json:"tls,omitempty"`
}

type LogMessage struct {
	Level       string              `json:"level"`
	Ts          float64             `json:"ts"`
	Logger      string              `json:"logger"`
	Msg         string              `json:"msg"`
	UserId      string              `json:"user_id"`
	Duration    float64             `json:"duration"`
	Size        int64               `json:"size"`
	Status      int                 `json:"status"`
	Request     LogRequest          `json:"request"`
	RespHeaders map[string][]string `json:"resp_headers"`
	/*
		Signal      string          `json:"signal,omitempty"`
		Cache       string          `json:"cache,omitempty"`
		Address     string          `json:"address,omitempty"`
		ExitCode    int             `json:"exit_code,omitempty"`
	*/
}

type LogInfo struct {
	TimeStamp      time.Time
	Duration       time.Duration
	Size           int64
	Status         int
	RemoteIp       string
	RemotePort     int
	Proto          string
	Method         string
	Host           string
	Uri            string
	UserAgent      string
	CfConnectingIp string
	CfIpcountry    string
	XForwardedFor  string
	TlsServerName  string
}

func main() {
	socketAddr := flag.String("s", "", "socket address")
	listenPort := flag.Int("p", 0, "listen port")
	dbPath := flag.String("d", "", "database path")
	flag.Parse()

	var conn net.Conn
	var err error
	if *socketAddr != "" {
		conn, err = getUnixConn(*socketAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *listenPort != 0 {
		conn, err = getUdpConn(*listenPort)
		if err != nil {
			log.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go readConn(conn, &wg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	conn.Close()
	wg.Wait()

	if *socketAddr != "" {
		os.Remove(*socketAddr)
	}

	fmt.Println()
	log.Println("Exiting")
}

func getUnixConn(socketAddr string) (net.Conn, error) {
	if _, err := os.Stat(socketAddr); err == nil {
		fmt.Println("Removing existing socket")
		if err := os.Remove(socketAddr); err != nil {
			fmt.Println("Error removing existing socket:", err)
			return nil, err
		}
	}

	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{
		Name: socketAddr,
		Net:  "unixgram",
	})
	if err != nil {
		fmt.Println("Error listening:", err)
		return nil, err
	}
	fmt.Println("Listening on", socketAddr)

	return conn, nil
}

func getUdpConn(listenPort int) (net.Conn, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: listenPort,
	})
	if err != nil {
		fmt.Println("Error listening:", err)
		return nil, err
	}
	fmt.Println("Listening on port", listenPort)

	return conn, nil
}

func readConn(conn net.Conn, wg *sync.WaitGroup) {
	for {
		buf := make([]byte, 1024*64)
		nr, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		_, info, err := getLogMessage(buf[:nr])
		if err != nil {
			fmt.Println("Error getting log message:", err)
			continue
		}
		infoStr, _ := json.Marshal(info)
		fmt.Println(string(infoStr) + "\n")
	}
	wg.Done()
}

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

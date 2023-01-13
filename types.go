package main

type Message struct {
	Level       string              `json:"level"`
	Ts          float64             `json:"ts"`
	Logger      string              `json:"logger"`
	Msg         string              `json:"msg"`
	UserId      string              `json:"user_id"`
	Duration    float64             `json:"duration"`
	Size        int64               `json:"size"`
	Status      int                 `json:"status"`
	Request     Request             `json:"request"`
	RespHeaders map[string][]string `json:"resp_headers"`
	/*
		Signal      string          `json:"signal,omitempty"`
		Cache       string          `json:"cache,omitempty"`
		Address     string          `json:"address,omitempty"`
		ExitCode    int             `json:"exit_code,omitempty"`
	*/
}

type Request struct {
	RemoteIp   string              `json:"remote_ip"`
	RemotePort string              `json:"remote_port"`
	Proto      string              `json:"proto"`
	Method     string              `json:"method"`
	Host       string              `json:"host"`
	Uri        string              `json:"uri"`
	Headers    map[string][]string `json:"headers"`
	Tls        *Tls                `json:"tls,omitempty"`
}

type Tls struct {
	Resumed     bool   `json:"resumed"`
	Version     int    `json:"version"`
	CipherSuite int    `json:"cipher_suite"`
	Proto       string `json:"proto"`
	ServerName  string `json:"server_name"`
}

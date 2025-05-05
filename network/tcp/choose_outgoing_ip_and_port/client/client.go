package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
)

func main() {
	var host string
	var port int
	var remoteHost string
	var remotePort int

	flag.StringVar(&host, "host", "127.0.0.1", "host for source")
	flag.IntVar(&port, "port", 10000, "port for source")
	flag.StringVar(&remoteHost, "remoteHost", "127.0.0.1", "host for remote")
	flag.IntVar(&remotePort, "remotePort", 8080, "port for remote")
	flag.Parse()	

	fmt.Printf("debug: local=%v:%v, remote=%v:%v\n", host, port, remoteHost, remotePort)

	addr := net.TCPAddr{IP: net.ParseIP(host), Port: port}

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				LocalAddr: &addr,
			}).DialContext,
		},
	}
	resp, err := client.Get(fmt.Sprintf("http://%v:%v", remoteHost, remotePort))
	if err != nil {
		panic(fmt.Errorf("when exeute request caused error: %v", err))
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("when parse data caused error: %v", err))
	}

	fmt.Printf("body: %v\n", string(data))
}
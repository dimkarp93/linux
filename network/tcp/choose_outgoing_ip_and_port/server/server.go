package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	var host string
	var port int

	flag.StringVar(&host, "host", "0.0.0.0", "host for bind")
	flag.IntVar(&port, "port", 8080, "port for bin")
	flag.Parse()

	fmt.Printf("debug: host=%v, port=%v\n", host, port)
	listener, err := net.Listen("tcp4", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		panic(fmt.Errorf("cannot listen port: %v", err))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(10) * time.Second)

		w.Write([]byte(fmt.Sprintf("echo from: remote=%v\n", r.RemoteAddr)))
	})

	http.Serve(listener, mux)
}
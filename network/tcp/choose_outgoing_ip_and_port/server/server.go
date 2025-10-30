package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"
)

func main() {
	var counts int
	var host string
	var port int

	flag.IntVar(&counts, "counts", 10, "count of ports for binding")
	flag.StringVar(&host, "host", "0.0.0.0", "host for bind")
	flag.IntVar(&port, "port", 8080, "port for bin")
	flag.Parse()

	var wg sync.WaitGroup

	for idx := range counts {
		wg.Add(1)
		go doBind(host, port + idx)
	}

	wg.Wait()
}

func doBind(host string, port int) {
	fmt.Printf("debug: host=%v, port=%v\n", host, port)
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if err != nil {
					fmt.Printf("Error setting SO_REUSEPORT: %v\n", err)
				}
			})
		},
	}
	listener, err := lc.Listen(context.Background(), "tcp4", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		panic(fmt.Errorf("cannot listen port: %v", err))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(10) * time.Second)

		w.Write([]byte(fmt.Sprintf("echo from: remote=%v, original ping=%v\n", r.RemoteAddr, r.URL.Query().Get("ping"))))
	})

	http.Serve(listener, mux)
}
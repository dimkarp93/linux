package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	var host string
	var port int
	var counts int
	var remoteHost string
	var remotePort int

	flag.StringVar(&host, "host", "127.0.0.1", "host for source")
	flag.IntVar(&port, "port", 10000, "port for source")
	flag.IntVar(&counts, "counts", 10, "count of connections")
	flag.StringVar(&remoteHost, "remoteHost", "127.0.0.1", "host for remote")
	flag.IntVar(&remotePort, "remotePort", 8080, "port for remote")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	addr := net.TCPAddr{IP: net.ParseIP(host), Port: port}

	var wg sync.WaitGroup

	for idx := range counts {
		wg.Add(1)
		go doCheck(&wg, addr, host, port, remoteHost, remotePort+idx)
	}

	wg.Wait()
}

func doCheck(wg *sync.WaitGroup, addr net.TCPAddr, host string, port int, remoteHost string, remotePort int) {
	ping := genPing()
	fmt.Printf("debug: local=%v:%v, remote=%v:%v, debug=%v\n", host, port, remoteHost, remotePort, ping)

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				LocalAddr: &addr,
				Control: func(network, address string, c syscall.RawConn) error {
					return c.Control(func(fd uintptr) {
						var err error
						
						err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
						if err != nil {
							fmt.Printf("Error setting SO_REUSEPORT: %v\n", err)
						}

						err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
						if err != nil {
							fmt.Printf("Error setting SO_REUSEPORT: %v\n", err)
						}
					})
				},
			}).DialContext,
		},
	}
	resp, err := client.Get(fmt.Sprintf("http://%v:%v?ping=%v", remoteHost, remotePort, ping))
	if err != nil {
		panic(fmt.Errorf("when exeute request caused error: %v", err))
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("when parse data caused error: %v", err))
	}

	fmt.Printf("remote=%v:%v body: %v\n", remoteHost, remotePort, process(data))
	wg.Done()
}

func genPing() string {
	return strconv.Itoa(rand.Int())
}

func process(data []byte) string {
	return strings.Trim(string(data), "\n")
}

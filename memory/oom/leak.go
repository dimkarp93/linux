package main

import (
	"flag"
	"math/rand"
	"sync"
	"time"
)

const LIMIT = 1024
const BATCH = 1024

func main() {
	var n int
	flag.IntVar(&n, "proc", 1, "process")
	
	var iter int
	flag.IntVar(&iter, "iter", 10, "iterations")

	var mem int
	flag.IntVar(&mem, "mem", 1024, "mem step")

	var sleep int
	flag.IntVar(&sleep, "sleep", 1000, "sleep")

	var maxiter int
	flag.IntVar(&maxiter, "max-iter", 0, "max-iter")

	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(n)
	
	for range n {
		go func () {
			var objs []byte

			idx := 0

			for {
				usingOld := maxiter > 0 && idx >= maxiter

				var data []byte

				if usingOld {
					data = objs[len(objs) - mem * BATCH:]
				} else {
					data = make([]byte, mem * BATCH)
					for i := range mem * BATCH {
						data[i] = byte(rand.Intn(LIMIT) & 127)	
					}
				}
			
				for i := range mem * BATCH {
					for range iter {
						data[i] *= data[i]
						data[i] &= 127
					}
				}

				if !usingOld {
					objs = append(objs, data...)
					idx++
				}

				time.Sleep(time.Duration(sleep) * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}
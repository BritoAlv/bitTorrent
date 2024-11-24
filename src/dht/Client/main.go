package main

import (
	common "Centralized/Common"
	"fmt"
	"net"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	var address = listener.Addr().String()
	var logger = common.NewLogger(fmt.Sprintf("Client%s.txt", address))
	name, err := register(address, logger)
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		speak(name, logger)
	}()
	go func() {
		defer wg.Done()
		listen(&listener, name, logger)
	}()
	wg.Wait()
}

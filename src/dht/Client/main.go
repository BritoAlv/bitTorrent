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
	wg.Add(1)
	done := make(chan struct{})

	go func() {
		defer wg.Done()
		speak(name, logger)
		close(done)
	}()

	go func() {
		defer wg.Done()
		select {
		case <-done:
			err := listener.Close()
			if err != nil {
				logger.WriteToFileError(fmt.Sprintf("Failed stopping listener,  was using address = %s, error was %s", listener.Addr().String(), err.Error()))
			}
			logger.WriteToFileOK(fmt.Sprintf("Succesfully Stopped Listener, was using address = %s", listener.Addr().String()))
			return
		default:
			listen(&listener, name, logger)
		}
	}()
	wg.Wait()
}

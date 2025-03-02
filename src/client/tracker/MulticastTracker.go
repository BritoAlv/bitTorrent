package tracker

import (
	"bittorrent/common"
	"bittorrent/dht/library/WithSocket"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type signal struct {
	mutex sync.Mutex
	state bool
}

func newSignal() *signal {
	return &signal{
		state: false,
		mutex: sync.Mutex{},
	}
}

func (s *signal) read() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.state
}

func (s *signal) write(value bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.state = value
}

type multicastTracker struct {
	MulticastUrl string
	ServerUrl    string
}

func (tracker *multicastTracker) Track(request common.TrackRequest) (common.TrackResponse, error) {
	const MAX_COUNT = 10
	count := 0

	multicastIp, multicastPort, err := getMulticastIpPort(tracker.MulticastUrl)
	if err != nil {
		return common.TrackResponse{}, err
	}

	var response common.TrackResponse
	for response, err = tracker.sendRequest(request); err != nil; {
		signal := newSignal()
		multicastChannel := make(chan [2]string)

		localIp, _ := WithSocket.GetIpFromInterface("eth0")
		listener, err := net.Listen("tcp", localIp+":")
		if err != nil {
			log.Println("Failed to start listener:", err)
			os.Exit(1)
		}
		localPort := strings.Split(listener.Addr().String(), ":")[1]

		go sendToMulticast(signal, multicastIp, multicastPort, localIp, localPort)
		go receiveFromMulticast(multicastChannel, listener)

		var serverIp, serverPort string
		for message := range multicastChannel {
			serverIp, serverPort = message[0], message[1]
			signal.write(true)
			break
		}

		tracker.ServerUrl = fmt.Sprintf("http://%s:%s/announce", serverIp, serverPort)

		count++
		if count > MAX_COUNT {
			return response, err
		}
	}
	return response, err
}

func (tracker *multicastTracker) sendRequest(request common.TrackRequest) (common.TrackResponse, error) {
	UrlSend, err := common.BuildHttpUrl(tracker.ServerUrl, request)

	if err != nil {
		return common.TrackResponse{}, fmt.Errorf("an error occurred while building the url to contact the tracker : %w", err)
	}
	// Send GET request
	httpResponse, err := http.Get(UrlSend)
	if err != nil {
		return common.TrackResponse{}, fmt.Errorf("an error occurred while contacting the tracker: %w", err)
	}

	bytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return common.TrackResponse{}, fmt.Errorf("an error ocurred while reading the body of the response from the tracker: %w", err)
	}

	response, err := common.DecodeTrackerResponse(bytes)
	if err != nil {
		return common.TrackResponse{}, fmt.Errorf("an error occurred while decoding the response from the tracker: %w", err)
	}

	log.Printf("Tracker's response: %v\n", response)
	if response.FailureReason != "" {
		return common.TrackResponse{}, errors.New(response.FailureReason)
	}
	return response, nil
}

func receiveFromMulticast(channel chan [2]string, listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}

		bytes, err := io.ReadAll(connection)
		if err != nil {
			log.Println("Failed to read from connection", err)
			connection.Close()
			continue
		}

		message := strings.Split(string(bytes), ";")
		if len(message) != 2 {
			log.Println("Invalid message")
			continue
		}
		serverIp, serverPort := message[0], message[1]
		channel <- [2]string{serverIp, serverPort}
		log.Printf("Received message: %v:%v \n", serverIp, serverPort)

		connection.Close()
		listener.Close()
		break
	}
}

func sendToMulticast(sig *signal, multicastIp string, multicastPort string, ip string, port string) {
	for {
		if sig.read() {
			break
		}

		// Resolve the UDP address
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", multicastIp, multicastPort))
		if err != nil {
			log.Println("Failed to resolve UDP address:", err)
		}

		// Create the UDP socket
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			log.Println("Failed to dial UDP:", err)
			os.Exit(1)
		}

		// Send message to the multicast group
		message := ip + ";" + port
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Println("Failed to send message:", err)
		}
		log.Println("Sent message to multicast group:", message)
		time.Sleep(time.Second * 3)
		conn.Close()
	}
}

func getMulticastIpPort(url string) (string, string, error) {
	urlParts := strings.Split(url, ":")
	if len(urlParts) < 3 {
		return "", "", errors.New("wrong multicast url")
	}
	return strings.TrimPrefix(urlParts[1], "//"), strings.Split(urlParts[2], "/")[0], nil
}

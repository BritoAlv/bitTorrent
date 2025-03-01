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

type CentralizedHttpTracker struct {
	MulticastUrl string
	ServerUrl    string
}

func (tracker *CentralizedHttpTracker) Track(request common.TrackRequest) (common.TrackResponse, error) {
	const LOCAL_PORT = "9090"
	MAX_COUNT := 10
	count := 0

	multicastIp, multicastPort, err := getMulticastIpPort(tracker.MulticastUrl)
	if err != nil {
		return common.TrackResponse{}, err
	}

	var response common.TrackResponse
	for response, err = tracker.sendRequest(tracker.ServerUrl, request); err != nil; {
		signal := newSignal()
		multicastChannel := make(chan [2]string)

		localIp, _ := WithSocket.GetIpFromInterface("eth0")

		go sendToMulticast(signal, multicastIp, multicastPort, localIp, LOCAL_PORT)
		go receiveFromMulticast(multicastChannel, localIp, LOCAL_PORT)

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

func (tracker *CentralizedHttpTracker) sendRequest(url string, request common.TrackRequest) (common.TrackResponse, error) {
	UrlSend, err := common.BuildHttpUrl(url, request)

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
	return response, nil
}

func receiveFromMulticast(channel chan [2]string, ip string, port string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", ip, port))
	if err != nil {
		log.Println("Failed to start listener:", err)
		os.Exit(1)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}

		bytes, err := io.ReadAll(connection)
		if err != nil {
			log.Println("Failed to read from connection", err)
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
		defer conn.Close()

		// Send message to the multicast group
		message := ip + ";" + port
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Println("Failed to send message:", err)
		}
		log.Println("Sent message to multicast group:", message)
		time.Sleep(time.Second * 3)
	}
}

func getMulticastIpPort(url string) (string, string, error) {
	urlParts := strings.Split(url, ":")
	if len(urlParts) < 3 {
		return "", "", errors.New("wrong multicast url")
	}
	return strings.TrimPrefix(urlParts[1], "//"), strings.Split(urlParts[2], "/")[0], nil
}

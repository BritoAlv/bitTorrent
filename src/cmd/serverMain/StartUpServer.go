package main

import (
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/MonitorHand"
	"bittorrent/dht/library/WithSocket"
	"flag"
	"strconv"
	"strings"
	"time"
)

func main() {
	iface := flag.String("iface", "tun0", "Network interface to use")
	ports := flag.String("ports", "12345,12346,12347,12348,12349", "Comma-separated list of available ports")
	flag.Parse()
	portList := strings.Split(*ports, ",")
	WithSocket.RegisterStartUp(*iface, "SocketServerClient", portList)
	randomId := Core.GenerateRandomBinaryId()
	randomIdStr := strconv.Itoa(int(randomId))
	socketServerClient := WithSocket.NewSocketServerClient(randomId)
	monitorHand := MonitorHand.NewMonitorHand[WithSocket.SocketContact]("Monitor" + randomIdStr)
	nodeSocket := Core.NewBruteChord(socketServerClient, socketServerClient, monitorHand, randomId)
	for {
		time.Sleep(5 * time.Second)
		println(nodeSocket.GetId())
	}
}

package main

import (
	"bittorrent/dht/library/WithSocket"
	"flag"
	"strings"
	"time"
)

func main() {
	iface := flag.String("iface", "tun0", "Network interface to use")
	portsUdp := flag.String("portsUdp", "12345,12346,12347,12348,12349", "Comma-separated list of available portsUdp")
	flag.Parse()
	portList := strings.Split(*portsUdp, ",")
	WithSocket.RegisterStartUp(*iface, "SocketServerClient", portList)
	nodeSocket := WithSocket.NewNodeSocket()
	for {
		time.Sleep(5 * time.Second)
		println(nodeSocket.GetId())
	}
}

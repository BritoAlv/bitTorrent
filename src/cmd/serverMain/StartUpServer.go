package main

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/MonitorHand"
	"bittorrent/dht/library/WithSocket"
	"strconv"
	"time"
)

func main() {
	Core.RegisterNotifications[WithSocket.SocketContact]()
	common.SetLogDirectoryPath("./SocketServerClient")
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

package WithSocket

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
)

var availablePortsUdp = []string{
	"12345",
	"12346",
	"12347",
	"12348",
	"12349",
	"12350",
	"12351",
	"12352",
	"12353",
	"12354",
	"12355",
	"12356",
	"12357",
	"12358",
	"12359",
	"12360",
	"12361",
}

var networkInterface = "eth0"

func SetNetworkInterface(iface string) {
	networkInterface = iface
}

func RegisterStartUp(iface string, name string, ports []string) {
	SetNetworkInterface(iface)
	availablePortsUdp = ports
	common.SetLogDirectoryPath("./" + name)
	Core.RegisterNotifications[SocketContact]()
}

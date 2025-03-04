package WithSocket

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/MonitorHand"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"time"
)

type Announcer struct {
	activeKnown     Core.SafeMap[Core.ChordHash, SocketContact]
	diagramListener net.PacketConn
	Contact         SocketContact
	logger          common.Logger
	monitor         Core.Monitor[SocketContact]
}

func NewAnnouncer(contact SocketContact) *Announcer {
	var announcer Announcer
	announcer.Contact = contact
	announcer.monitor = MonitorHand.NewMonitorHand[SocketContact]("MonitorAnnouncer" + strconv.Itoa(int(contact.GetNodeId())) + ".txt")
	announcer.activeKnown = *Core.NewSafeMap[Core.ChordHash, SocketContact](make(map[Core.ChordHash]SocketContact))
	announcer.logger = *common.NewLogger("Announcer" + strconv.Itoa(int(contact.GetNodeId())) + ".txt")
	_, broadIP := GetIpFromInterface(networkInterface)
	randomPort := availablePortsUdp[rand.Int()%len(availablePortsUdp)]
	add, err := net.ResolveUDPAddr("udp", broadIP+":"+randomPort)
	if err != nil {
		announcer.logger.WriteToFileError("Error resolving UDP address: %v", err)
	}
	announcer.addContact(contact)
	diagramListener, err := net.ListenUDP("udp", add)
	if err != nil {
		announcer.logger.WriteToFileError("Error Listening %v", err)
		panic(err)
	}
	announcer.logger.WriteToFileOK("Listening UDP on %v", add.String())
	announcer.diagramListener = diagramListener
	go announcer.listenAnnounces()
	go announcer.sendAnnounces()
	return &announcer
}

func (a *Announcer) addContact(contact SocketContact) {
	a.activeKnown.Set(contact.GetNodeId(), contact)
	a.monitor.AddContact(contact, time.Now())
	for _, knownContact := range a.activeKnown.GetValues() {
		if !a.monitor.CheckAlive(knownContact, 3) {
			a.activeKnown.Delete(knownContact.GetNodeId())
			a.monitor.DeleteContact(knownContact)
		}
	}
}

func (a *Announcer) GetContacts() []SocketContact {
	return a.activeKnown.GetValues()
}

func (a *Announcer) listenAnnounces() {
	buf := make([]byte, 4096)
	for {
		n, addr, err := a.diagramListener.ReadFrom(buf)
		if err != nil {
			a.logger.WriteToFileError("Error reading: %v", err)
			continue
		}
		a.logger.WriteToFileOK("Received data from %v", addr)
		reader := bytes.NewReader(buf[:n])
		gobDecoder := gob.NewDecoder(reader)
		var receivedContact SocketContact
		err = gobDecoder.Decode(&receivedContact)
		if err != nil {
			a.logger.WriteToFileError("Error Decoding %v", err)
			continue
		}
		fmt.Printf("Received Contact with nodeID = %v, and address = %v \n", receivedContact.NodeId, receivedContact.Addr)
		go a.addContact(receivedContact)
		a.logger.WriteToFileOK("Received %v from %v", receivedContact, addr)
	}
}

func (a *Announcer) sendAnnouncesLogic() {
	_, broadcastAddr := GetIpFromInterface(networkInterface)
	for _, port := range availablePortsUdp {
		conn, err := net.Dial("udp", broadcastAddr+":"+port)
		if err != nil {
			a.logger.WriteToFileError("Error dialing: %v", err)
			continue
		}
		message := a.Contact
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		err = encoder.Encode(message)
		if err != nil {
			a.logger.WriteToFileError("Error encoding: %v", err)
			continue
		}
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			a.logger.WriteToFileError("Error writing: %v", err)
			continue
		}
		a.logger.WriteToFileOK("Broadcast Message : %v,  Sent to %v from %v", message, broadcastAddr, conn.LocalAddr())
		err = conn.Close()
		if err != nil {
			a.logger.WriteToFileError("Error closing: %v", err)
			continue
		}
	}
}

func (a *Announcer) sendAnnounces() {
	for {
		a.logger.WriteToFileOK("Known Contacts are: %v", a.GetContacts())
		time.Sleep(3 * time.Second)
		a.sendAnnouncesLogic()
	}
}

package WithSocket

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"sync"
	"time"
)

type Announcer struct {
	activeKnown     map[Core.ChordHash]SocketContact
	diagramListener net.PacketConn
	Contact         SocketContact
	lock            sync.Mutex
	logger          common.Logger
}

func NewAnnouncer(contact SocketContact) *Announcer {
	var announcer Announcer
	announcer.Contact = contact
	announcer.lock = sync.Mutex{}
	announcer.activeKnown = make(map[Core.ChordHash]SocketContact)
	announcer.logger = *common.NewLogger("Announcer" + strconv.Itoa(int(contact.GetNodeId())) + ".txt")
	_, broadIP := getIPFromInterface()
	add, err := net.ResolveUDPAddr("udp", broadIP+":"+usedPorts[rand.Int()%len(usedPorts)])
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
	a.lock.Lock()
	defer a.lock.Unlock()
	a.activeKnown[contact.GetNodeId()] = contact
}

func (a *Announcer) GetContacts() []SocketContact {
	a.lock.Lock()
	defer a.lock.Unlock()
	result := make([]SocketContact, 0)
	for _, contact := range a.activeKnown {
		result = append(result, contact)
	}
	return result
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
	_, broadcastAddr := getIPFromInterface()
	for _, port := range usedPorts {
		conn, err := net.Dial("udp", broadcastAddr+":"+port)
		if err != nil {
			a.logger.WriteToFileError("Error dialing: %v", err)
			return
		}
		message := a.Contact
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		err = encoder.Encode(message)
		if err != nil {
			a.logger.WriteToFileError("Error encoding: %v", err)
			return
		}
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			a.logger.WriteToFileError("Error writing: %v", err)
			return
		}
		a.logger.WriteToFileOK("Broadcast Message : %v,  Sent to %v from %v", message, broadcastAddr, conn.LocalAddr())
		err = conn.Close()
		if err != nil {
			a.logger.WriteToFileError("Error closing: %v", err)
			return
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

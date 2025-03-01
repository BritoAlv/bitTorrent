package Core

import (
	"encoding/gob"
	"net"
)

const NumberBits = 8
const WaitingTime = 1
const Attempts = 5

type ChordHash = int64
type Store = map[ChordHash][]byte

func RegisterNotifications[T Contact]() {
	gob.Register(&areYouAliveNotification[T]{})
	gob.Register(&imAliveNotification[T]{})
	gob.Register(&areYouMyPredecessor[T]{})
	gob.Register(&imYourPredecessor[T]{})
	gob.Register(&getRequest[T]{})
	gob.Register(&receivedGetRequest[T]{})
	gob.Register(&putRequest[T]{})
	gob.Register(&receivedPutRequest[T]{})
	gob.Register(&receiveDataReplicate[T]{})
	gob.Register(&confirmReplication[T]{})
	gob.Register(&TellMeYourState[T]{})
	gob.Register(&TellMeYourStateResponse[T]{})
	gob.Register(&net.TCPAddr{})
}

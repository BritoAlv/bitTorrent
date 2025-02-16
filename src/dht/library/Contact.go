package library

/*
This should contain enough information so that a client (Chord Node with the need to send a request)
can communicate with a server (Chord Node waiting to receive a request).
*/

type Contact interface {
	getNodeId() [NumberBits]uint8
}

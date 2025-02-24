package library

/*
The Server will be doing the function of listening, so If someone wants to communicate with a Node it should to with
the server, The server is responsible for:

The Server forwarding the notifications to the Node.

Server -> Node is a one to one correspondence.

The T Contact used determines the type of communication that the server should use.

Because a priori, the server does not know the ID and the channel of the Node that he will be responsible it should have
a method for handling that.

*/

type Server[T Contact] interface {
	GetContact() T
	SetData(channel chan Notification[T], NodeId ChordHash)
}

package messenger

import "io"

const _HANDSHAKE_MESSAGE = -1
const _CHOKE_MESSAGE = 0
const _UNCHOKE_MESSAGE = 1
const _INTERESTED_MESSAGE = 2
const _NOT_INTERESTED_MESSAGE = 3
const _HAVE_MESSAGE = 4
const _BITFIELD_MESSAGE = 5
const _REQUEST_MESSAGE = 6
const _PIECE_MESSAGE = 7
const _CANCEL_MESSAGE = 8

type Messenger interface {
	Write(writer io.Writer, message interface{}) error
	Read(io.Reader) (interface{}, error)
}

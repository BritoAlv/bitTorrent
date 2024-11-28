package messenger

import (
	"bittorrent/common"
	"errors"
	"io"
	"strconv"
	"strings"
)

type stringifiedMessenger struct {
}

func New() Messenger {
	return stringifiedMessenger{}
}

//** Write implementation

func (manager stringifiedMessenger) Write(writer io.Writer, message interface{}) error {
	var bytes []byte
	switch castedMessage := message.(type) {
	case HandshakeMessage:
		bytes = encodeHandshakeMessage(castedMessage)
	case ChokeMessage:
		bytes = encodePayloadLessMessage(_CHOKE_MESSAGE)
	case UnchokeMessage:
		bytes = encodePayloadLessMessage(_UNCHOKE_MESSAGE)
	case InterestedMessage:
		bytes = encodePayloadLessMessage(_INTERESTED_MESSAGE)
	case NotInterestedMessage:
		bytes = encodePayloadLessMessage(_NOT_INTERESTED_MESSAGE)
	case HaveMessage:
		bytes = encodeHaveMessage(castedMessage)
	case BitfieldMessage:
		bytes = encodeBitfieldMessage(castedMessage)
	case RequestMessage:
		bytes = encodeRequestMessage(castedMessage)
	case PieceMessage:
		bytes = encodePieceMessage(castedMessage)
	case CancelMessage:
		bytes = encodeCancelMessage(castedMessage)
	default:
		return errors.New("invalid message type")
	}

	err := common.ReliableWrite(writer, bytes)
	if err != nil {
		return err
	}
	return nil
}

func encodeHandshakeMessage(message HandshakeMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_HANDSHAKE_MESSAGE) + ";" + message.Id + ";")
	messageBytes = append(messageBytes, message.Infohash[:]...)
	return append(getLength(messageBytes), messageBytes...)
}

func encodePayloadLessMessage(messageType int) []byte {
	message := []byte(strconv.Itoa(messageType) + ";")
	return append(getLength(message), message...)
}

func encodeHaveMessage(message HaveMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_HAVE_MESSAGE) + ";" + strconv.Itoa(message.Index) + ";")
	return append(getLength(messageBytes), []byte(messageBytes)...)
}

func encodeBitfieldMessage(message BitfieldMessage) []byte {
	messageStr := strconv.Itoa(_BITFIELD_MESSAGE) + ";"

	for _, bit := range message.Bitfield {
		if bit {
			messageStr += "1;"
		} else {
			messageStr += "0;"
		}
	}

	messageBytes := []byte(messageStr)

	return append(getLength(messageBytes), messageBytes...)
}

func encodeRequestMessage(message RequestMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_REQUEST_MESSAGE) + ";" + strconv.Itoa(message.Index) + ";" + strconv.Itoa(message.Offset) + ";" + strconv.Itoa(message.Length) + ";")
	return append(getLength(messageBytes), messageBytes...)
}

func encodePieceMessage(message PieceMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_PIECE_MESSAGE) + ";" + strconv.Itoa(message.Index) + ";" + strconv.Itoa(message.Offset) + ";")

	messageBytes = append(messageBytes, message.Bytes...)

	return append(getLength(messageBytes), messageBytes...)
}

func encodeCancelMessage(message CancelMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_CANCEL_MESSAGE) + ";" + strconv.Itoa(message.Index) + ";" + strconv.Itoa(message.Offset) + ";" + strconv.Itoa(message.Length) + ";")
	return append(getLength(messageBytes), messageBytes...)
}

func getLength(message []byte) []byte {
	messageLength := ";" + strconv.Itoa(len(message)) + ";"
	metaLength := byte(len(messageLength))
	return append([]byte{metaLength}, []byte(messageLength)...)
}

//** Read implementation

func (manager stringifiedMessenger) Read(reader io.Reader) (interface{}, error) {
	metaLengthBytes, err := common.ReliableRead(reader, 1)
	if err != nil {
		return nil, err
	}

	metaLength := int(metaLengthBytes[0])
	messageLengthBytes, err := common.ReliableRead(reader, metaLength)
	if err != nil {
		return nil, err
	}

	messageLength, err := strconv.Atoi(string(messageLengthBytes[1 : metaLength-1]))
	if err != nil {
		return nil, err
	}

	messageBytes, err := common.ReliableRead(reader, messageLength)
	if err != nil {
		return nil, err
	}

	messageStr := string(messageBytes)
	messageSplits := strings.Split(messageStr, ";")
	messageType, err := strconv.Atoi(messageSplits[0])
	if err != nil {
		return nil, err
	}

	// Remove last empty-split if it's present
	if messageSplits[len(messageSplits)-1] == "" {
		messageSplits = messageSplits[:len(messageSplits)-1]
	}
	// Ignore first type-split for now on
	messageSplits = messageSplits[1:]

	switch messageType {
	case _HANDSHAKE_MESSAGE:
		return decodeHandshakeMessage(messageSplits)
	case _CHOKE_MESSAGE:
		return ChokeMessage{}, nil
	case _UNCHOKE_MESSAGE:
		return UnchokeMessage{}, nil
	case _INTERESTED_MESSAGE:
		return InterestedMessage{}, nil
	case _NOT_INTERESTED_MESSAGE:
		return NotInterestedMessage{}, nil
	case _HAVE_MESSAGE:
		return decodeHaveMessage(messageSplits)
	case _BITFIELD_MESSAGE:
		return decodeBitfieldMessage(messageSplits)
	case _REQUEST_MESSAGE:
		return decodeRequestMessage(messageSplits)
	case _PIECE_MESSAGE:
		return decodePieceMessage(messageSplits)
	case _CANCEL_MESSAGE:
		return decodeCancelMessage(messageSplits)
	default:
		return nil, errors.New("invalid message type")
	}
}

func decodeHandshakeMessage(messageSplits []string) (HandshakeMessage, error) {
	if len(messageSplits) != 2 {
		return HandshakeMessage{}, errors.New("invalid handshake-message payload")
	}

	id := string(messageSplits[0])

	infohashSlice := []byte(messageSplits[1])
	if len(infohashSlice) != 20 {
		return HandshakeMessage{}, errors.New("invalid handshake-message payload")
	}
	infohash := [20]byte(infohashSlice)

	return HandshakeMessage{
		Infohash: infohash,
		Id:       id,
	}, nil
}

func decodeHaveMessage(messageSplits []string) (HaveMessage, error) {
	if len(messageSplits) != 1 {
		return HaveMessage{}, errors.New("invalid have-message payload")
	}

	index, err := strconv.Atoi(messageSplits[0])
	if err != nil {
		return HaveMessage{}, errors.New("invalid have-message payload")
	}

	return HaveMessage{Index: index}, nil
}

func decodeBitfieldMessage(messageSplits []string) (BitfieldMessage, error) {
	bitfield := []bool{}
	for _, bit := range messageSplits {
		if bit == "1" {
			bitfield = append(bitfield, true)
		} else if bit == "0" {
			bitfield = append(bitfield, false)
		} else {
			return BitfieldMessage{}, errors.New("invalid bitfield-message payload")
		}
	}

	return BitfieldMessage{Bitfield: bitfield}, nil
}

func decodeRequestMessage(messageSplits []string) (RequestMessage, error) {
	if len(messageSplits) != 3 {
		return RequestMessage{}, errors.New("invalid request-message payload")
	}

	index, err := strconv.Atoi(messageSplits[0])
	if err != nil {
		return RequestMessage{}, errors.New("invalid request-message payload")
	}

	offset, err := strconv.Atoi(messageSplits[1])
	if err != nil {
		return RequestMessage{}, errors.New("invalid request-message payload")
	}

	length, err := strconv.Atoi(messageSplits[2])
	if err != nil {
		return RequestMessage{}, errors.New("invalid request-message payload")
	}

	return RequestMessage{
		Index:  index,
		Offset: offset,
		Length: length,
	}, nil
}

func decodePieceMessage(messageSplits []string) (PieceMessage, error) {
	if len(messageSplits) != 3 {
		return PieceMessage{}, errors.New("invalid piece-message payload")
	}

	index, err := strconv.Atoi(messageSplits[0])
	if err != nil {
		return PieceMessage{}, errors.New("invalid piece-message payload")
	}

	offset, err := strconv.Atoi(messageSplits[1])
	if err != nil {
		return PieceMessage{}, errors.New("invalid piece-message payload")
	}

	bytes := []byte(messageSplits[2])

	return PieceMessage{
		Index:  index,
		Offset: offset,
		Bytes:  bytes,
	}, nil
}

func decodeCancelMessage(messageSplits []string) (CancelMessage, error) {
	if len(messageSplits) != 3 {
		return CancelMessage{}, errors.New("invalid cancel-message payload")
	}

	index, err := strconv.Atoi(messageSplits[0])
	if err != nil {
		return CancelMessage{}, errors.New("invalid cancel-message payload")
	}

	offset, err := strconv.Atoi(messageSplits[1])
	if err != nil {
		return CancelMessage{}, errors.New("invalid cancel-message payload")
	}

	length, err := strconv.Atoi(messageSplits[2])
	if err != nil {
		return CancelMessage{}, errors.New("invalid cancel-message payload")
	}

	return CancelMessage{
		RequestMessage{
			Index:  index,
			Offset: offset,
			Length: length,
		},
	}, nil
}

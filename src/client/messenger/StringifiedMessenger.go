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
	message := []byte(strconv.Itoa(messageType))
	return append(getLength(message), message...)
}

func encodeHaveMessage(message HaveMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_HAVE_MESSAGE) + ";" + strconv.Itoa(message.Index))
	return append(getLength(messageBytes), []byte(messageBytes)...)
}

func encodeBitfieldMessage(message BitfieldMessage) []byte {
	messageStr := strconv.Itoa(_BITFIELD_MESSAGE) + ";"

	for _, bit := range message.Bitfield {
		if bit {
			messageStr += "1"
		} else {
			messageStr += "0"
		}
	}

	messageBytes := []byte(messageStr)

	return append(getLength(messageBytes), messageBytes...)
}

func encodeRequestMessage(message RequestMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_REQUEST_MESSAGE) + ";" + strconv.Itoa(message.Index) + ";" + strconv.Itoa(message.Offset) + ";" + strconv.Itoa(message.Length))
	return append(getLength(messageBytes), messageBytes...)
}

func encodePieceMessage(message PieceMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_PIECE_MESSAGE) + ";" + strconv.Itoa(message.Index) + ";" + strconv.Itoa(message.Offset) + ";")

	messageBytes = append(messageBytes, message.Bytes...)

	return append(getLength(messageBytes), messageBytes...)
}

func encodeCancelMessage(message CancelMessage) []byte {
	messageBytes := []byte(strconv.Itoa(_CANCEL_MESSAGE) + ";" + strconv.Itoa(message.Index) + ";" + strconv.Itoa(message.Offset) + ";" + strconv.Itoa(message.Length))
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
	messageSplits := strings.SplitN(messageStr, ";", 2)
	messageType, err := strconv.Atoi(messageSplits[0])

	if err != nil {
		return nil, err
	}

	switch messageType {
	case _HANDSHAKE_MESSAGE:
		return decodeHandshakeMessage(messageStr)
	case _CHOKE_MESSAGE:
		return ChokeMessage{}, nil
	case _UNCHOKE_MESSAGE:
		return UnchokeMessage{}, nil
	case _INTERESTED_MESSAGE:
		return InterestedMessage{}, nil
	case _NOT_INTERESTED_MESSAGE:
		return NotInterestedMessage{}, nil
	case _HAVE_MESSAGE:
		return decodeHaveMessage(messageStr)
	case _BITFIELD_MESSAGE:
		return decodeBitfieldMessage(messageStr)
	case _REQUEST_MESSAGE:
		return decodeRequestMessage(messageStr)
	case _PIECE_MESSAGE:
		return decodePieceMessage(messageStr)
	case _CANCEL_MESSAGE:
		return decodeCancelMessage(messageStr)
	default:
		return nil, errors.New("invalid message type")
	}
}

func decodeHandshakeMessage(messageStr string) (HandshakeMessage, error) {
	messageSplits := strings.SplitN(messageStr, ";", 3)
	handshakeSplits := messageSplits[1:]

	if len(handshakeSplits) != 2 {
		return HandshakeMessage{}, errors.New("invalid handshake-message payload")
	}

	id := string(handshakeSplits[0])

	infohashSlice := []byte(handshakeSplits[1])
	if len(infohashSlice) != 20 {
		return HandshakeMessage{}, errors.New("invalid handshake-message payload")
	}
	infohash := [20]byte(infohashSlice)

	return HandshakeMessage{
		Infohash: infohash,
		Id:       id,
	}, nil
}

func decodeHaveMessage(messageStr string) (HaveMessage, error) {
	messageSplits := strings.SplitN(messageStr, ";", 2)
	haveSplits := messageSplits[1:]

	if len(haveSplits) != 1 {
		return HaveMessage{}, errors.New("invalid have-message payload")
	}

	index, err := strconv.Atoi(haveSplits[0])
	if err != nil {
		return HaveMessage{}, errors.New("invalid have-message payload")
	}

	return HaveMessage{Index: index}, nil
}

func decodeBitfieldMessage(messageStr string) (BitfieldMessage, error) {
	messageSplits := strings.SplitN(messageStr, ";", 2)
	bitfieldSplits := messageSplits[1:]

	if len(bitfieldSplits) != 1 {
		return BitfieldMessage{}, errors.New("invalid bitfield-message payload")
	}

	bitfield := []bool{}
	for _, bit := range bitfieldSplits[0] {
		if bit == 49 {
			bitfield = append(bitfield, true)
		} else if bit == 48 {
			bitfield = append(bitfield, false)
		} else {
			return BitfieldMessage{}, errors.New("invalid bitfield-message payload")
		}
	}

	return BitfieldMessage{Bitfield: bitfield}, nil
}

func decodeRequestMessage(messageStr string) (RequestMessage, error) {
	messageSplits := strings.SplitN(messageStr, ";", 4)
	requestSplits := messageSplits[1:]

	if len(requestSplits) != 3 {
		return RequestMessage{}, errors.New("invalid request-message payload")
	}

	index, err := strconv.Atoi(requestSplits[0])
	if err != nil {
		return RequestMessage{}, errors.New("invalid request-message payload")
	}

	offset, err := strconv.Atoi(requestSplits[1])
	if err != nil {
		return RequestMessage{}, errors.New("invalid request-message payload")
	}

	length, err := strconv.Atoi(requestSplits[2])
	if err != nil {
		return RequestMessage{}, errors.New("invalid request-message payload")
	}

	return RequestMessage{
		Index:  index,
		Offset: offset,
		Length: length,
	}, nil
}

func decodePieceMessage(messageStr string) (PieceMessage, error) {
	messageSplits := strings.SplitN(messageStr, ";", 4)
	pieceSplits := messageSplits[1:]

	if len(pieceSplits) != 3 {
		return PieceMessage{}, errors.New("invalid piece-message payload")
	}

	index, err := strconv.Atoi(pieceSplits[0])
	if err != nil {
		return PieceMessage{}, errors.New("invalid piece-message payload")
	}

	offset, err := strconv.Atoi(pieceSplits[1])
	if err != nil {
		return PieceMessage{}, errors.New("invalid piece-message payload")
	}

	bytes := []byte(pieceSplits[2])

	return PieceMessage{
		Index:  index,
		Offset: offset,
		Bytes:  bytes,
	}, nil
}

func decodeCancelMessage(messageStr string) (CancelMessage, error) {
	messageSplits := strings.SplitN(messageStr, ";", 4)
	cancelSplits := messageSplits[1:]

	if len(cancelSplits) != 3 {
		return CancelMessage{}, errors.New("invalid cancel-message payload")
	}

	index, err := strconv.Atoi(cancelSplits[0])
	if err != nil {
		return CancelMessage{}, errors.New("invalid cancel-message payload")
	}

	offset, err := strconv.Atoi(cancelSplits[1])
	if err != nil {
		return CancelMessage{}, errors.New("invalid cancel-message payload")
	}

	length, err := strconv.Atoi(cancelSplits[2])
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

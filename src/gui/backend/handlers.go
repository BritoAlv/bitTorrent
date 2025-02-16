package backend

import (
	"bittorrent/client/peer"
	"bittorrent/common"
	"bittorrent/torrent"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

type DownloadHandler struct{ Peers map[string]*peer.Peer }
type UpdateHandler struct{ Peers map[string]*peer.Peer }
type KillHandler struct{ Peers map[string]*peer.Peer }

func (handler *DownloadHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := getBody(request.Body)
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: err.Error(),
		})
		return
	}

	var downloadRequest DownloadRequest
	err = json.Unmarshal(body, &downloadRequest)
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: err.Error(),
		})
		return
	}

	_peer, err := startPeer(downloadRequest.Id, downloadRequest.TorrentPath, downloadRequest.DownloadPath, downloadRequest.IpAddress, downloadRequest.EncryptionLevel)
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: err.Error(),
		})
		return
	}

	go _peer.Torrent(nil)
	handler.Peers[downloadRequest.Id] = _peer

	respond(responseWriter, BooleanResponse{
		Successful:   true,
		ErrorMessage: "",
	})
}

func (handler *UpdateHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	queryMap := request.URL.Query()
	id := queryMap.Get("id")

	if id == "" {
		respond(responseWriter, UpdateResponse{
			BooleanResponse: BooleanResponse{
				Successful:   false,
				ErrorMessage: "Expecting 'id' query",
			},
			Progress: 0,
			Peers:    0,
		})
		return
	}

	_peer := handler.Peers[id]
	progress, peers := _peer.Status()

	respond(responseWriter, UpdateResponse{
		BooleanResponse: BooleanResponse{
			Successful:   true,
			ErrorMessage: "",
		},
		Progress: progress,
		Peers:    peers,
	})
}

func (handler *KillHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	queryMap := request.URL.Query()
	id := queryMap.Get("id")

	if id == "" {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: "Expecting 'id' query",
		})
		return
	}

	respond(responseWriter, BooleanResponse{
		Successful:   true,
		ErrorMessage: "Killed" + id,
	})
}

func getBody(body io.ReadCloser) ([]byte, error) {
	buffer := make([]byte, 2048)
	totalRead, err := body.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	return buffer[:totalRead], nil
}

func respond(responseWriter io.Writer, response interface{}) {
	responseBytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: err.Error(),
		})
		return
	}

	err = common.ReliableWrite(responseWriter, responseBytes)
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: err.Error(),
		})
		return
	}
}

func startPeer(id string, torrentPath string, downloadPath string, ip string, encryptionLevel bool) (*peer.Peer, error) {
	torrent, err := torrent.ParseTorrentFile(torrentPath)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", ip+":")
	if err != nil {
		return nil, err
	}

	_peer, err := peer.New(id, listener, torrent, downloadPath, encryptionLevel)
	if err != nil {
		return nil, err
	}

	return &_peer, nil
}

// func getFile(path string) (*os.File, error) {
// 	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, common.ALL_RW_PERMISSION)
// 	if err != nil {
// 		return nil, errors.New("error accessing the file: " + err.Error())
// 	}

// 	return file, nil
// }

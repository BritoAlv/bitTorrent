package main

import (
	"bittorrent/client/peer"
	"bittorrent/dht/library/WithSocket"
	"bittorrent/torrent"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow any origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup API
	handler := newApiHandler()
	router.POST("/download", handler.download)
	router.GET("/update/:id", handler.update)
	router.GET("/kill/:id", handler.kill)

	// Host files
	router.Static("/home", "./static")

	err := router.Run(":8080")
	if err != nil {
		log.Println("Server could not be started")
	}
}

// ** Contracts
type downloadRequest struct {
	Id              string
	TorrentPath     string
	DownloadPath    string
	EncryptionLevel bool
}

type booleanResponse struct {
	Successful   bool
	ErrorMessage string
}

type updateResponse struct {
	booleanResponse
	Progress float32
	Peers    int
}

// ** Handler definition and implementation
type apiHandler struct {
	Peers map[string]*peer.Peer
}

func newApiHandler() *apiHandler {
	return &apiHandler{
		Peers: make(map[string]*peer.Peer),
	}
}

func (handler *apiHandler) download(ctx *gin.Context) {
	var request downloadRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.IndentedJSON(http.StatusOK, booleanResponse{
			Successful:   false,
			ErrorMessage: err.Error(),
		})
		return
	}

	ipAddress, _ := WithSocket.GetIpFromInterface("eth0")
	if ipAddress == "" {
		ipAddress = "localhost"
	}

	_peer, err := startPeer(request.Id, request.TorrentPath, request.DownloadPath, ipAddress, request.EncryptionLevel)
	if err != nil {
		fmt.Println("API: Error: " + err.Error())
		ctx.IndentedJSON(http.StatusOK, booleanResponse{
			Successful:   false,
			ErrorMessage: err.Error(),
		})
		return
	}

	go _peer.Torrent(nil)
	handler.Peers[request.Id] = _peer

	ctx.IndentedJSON(http.StatusOK, booleanResponse{
		Successful:   true,
		ErrorMessage: "",
	})
}

func (handler *apiHandler) update(ctx *gin.Context) {
	id, found := ctx.Params.Get("id")
	if !found {
		ctx.IndentedJSON(http.StatusOK, updateResponse{
			booleanResponse: booleanResponse{
				Successful:   false,
				ErrorMessage: "Expecting 'id' query",
			},
			Progress: 0,
			Peers:    0,
		})
		return
	}

	_peer, contained := handler.Peers[id]
	if !contained {
		ctx.IndentedJSON(http.StatusOK, updateResponse{
			booleanResponse: booleanResponse{
				Successful:   false,
				ErrorMessage: "Non existing 'id'",
			},
			Progress: 0,
			Peers:    0,
		})
		return
	}

	progress, peers := _peer.Status()

	ctx.IndentedJSON(http.StatusOK, updateResponse{
		booleanResponse: booleanResponse{
			Successful:   true,
			ErrorMessage: "",
		},
		Progress: progress,
		Peers:    peers,
	})
}

func (handler *apiHandler) kill(ctx *gin.Context) {
	id, found := ctx.Params.Get("id")
	if !found {
		ctx.IndentedJSON(http.StatusOK, booleanResponse{
			Successful:   false,
			ErrorMessage: "Expecting 'id' query",
		})
		return
	}

	_peer, contained := handler.Peers[id]
	if !contained {
		ctx.IndentedJSON(http.StatusOK, booleanResponse{
			Successful:   false,
			ErrorMessage: "Non existing 'id'",
		})
		return
	}

	_peer.NotificationChannel <- peer.KillNotification{}
	delete(handler.Peers, id)

	ctx.IndentedJSON(http.StatusOK, booleanResponse{
		Successful:   true,
		ErrorMessage: "",
	})
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

package main

import (
	"bittorrent/client/peer"
	"bittorrent/gui/backend"
	"net/http"
)

func main() {
	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	requestMultiplexer := http.NewServeMux()
	peers := make(map[string]*peer.Peer)

	// Register the routes and handlers
	requestMultiplexer.Handle("/download", &backend.DownloadHandler{Peers: peers})
	requestMultiplexer.Handle("/update", &backend.UpdateHandler{Peers: peers})
	requestMultiplexer.Handle("/kill", &backend.KillHandler{Peers: peers})

	// Run the server
	http.ListenAndServe("127.0.0.1:9090", requestMultiplexer)
}

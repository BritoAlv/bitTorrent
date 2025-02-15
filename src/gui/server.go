package main

import (
	"bittorrent/gui/backend"
	"net/http"
)

func main() {
	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	requestMultiplexer := http.NewServeMux()

	// Register the routes and handlers
	requestMultiplexer.Handle("/download", &backend.DownloadHandler{})
	requestMultiplexer.Handle("/update", &backend.UpdateHandler{})
	requestMultiplexer.Handle("/kill", &backend.KillHandler{})

	// Run the server
	http.ListenAndServe("127.0.0.1:9090", requestMultiplexer)
}

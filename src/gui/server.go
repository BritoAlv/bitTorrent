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

	// Wrap the multiplexer with a CORS handler
	handlerWithCORS := corsMiddleware(requestMultiplexer)

	// Run the server
	http.ListenAndServe("127.0.0.1:9595", handlerWithCORS)
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Allowed headers

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content
			return
		}

		// Pass request to the next handler
		next.ServeHTTP(w, r)
	})
}

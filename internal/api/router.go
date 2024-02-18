package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"mucutGo/internal/yt"
	"net/http"
	"strings"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/fetch-video-info", VideoInfoHandler).Methods("POST")

	router.HandleFunc("/download-mp3-bulk", Mp3DownloadBulkHandler).Methods("POST")

	router.HandleFunc("/download-mp3", Mp3DownloadHandler).Methods("POST")

	router.HandleFunc("/update-mp3-metadata", Mp3UpdateMetadataHandler).Methods("POST")

	router.HandleFunc("/update-mp3-metadata-bulk", Mp3UpdateMetadataBulkHandler).Methods("POST")

	router.HandleFunc("/pack-siq", CreateSiqHandler).Methods("GET")

	router.HandleFunc("/ws/video-info", VideoInfoWebSocketHandler)

	return router
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allows connections from any origin
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections = make(map[*websocket.Conn]bool)

// VideoInfoWebSocketHandler replaces VideoInfoHandler for WebSocket usage
func VideoInfoWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()
	connections[conn] = true // Track the new connection

	log.Println("WebSocket connection established")

	for {
		if _, ok := connections[conn]; !ok {
			// If the connection is not tracked, ignore processing
			continue
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("Received message: %s\n", message)

		// Process the message as you would process the body in VideoInfoHandler
		titles := strings.Split(string(message), "\n")
		metadata, err := yt.FetchVideoMetadataFromText(titles)
		if err != nil {
			log.Printf("Error prefetching videos: %v", err)
			// Optionally, send an error message back to the client
			continue
		}

		responseData, err := json.Marshal(metadata)
		if err != nil {
			log.Printf("Error marshaling metadata: %v", err)
			// Optionally, send an error message back to the client
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, responseData); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

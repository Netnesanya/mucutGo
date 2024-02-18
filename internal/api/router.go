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
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Assuming the message is a newline-delimited list of titles.
		titles := strings.Split(string(message), "\n")
		for _, title := range titles {
			// Fetch metadata for each title individually.
			metadata, err := yt.FetchVideoMetadataFromText([]string{title})
			if err != nil {
				log.Printf("Error fetching video metadata for title %s: %v", title, err)
				// Optionally, send an error message back to the client.
				continue
			}

			// Send back the metadata as soon as it's fetched.
			responseData, err := json.Marshal(metadata)
			if err != nil {
				log.Printf("Error marshaling metadata for title %s: %v", title, err)
				// Optionally, send an error message back to the client.
				continue
			}

			if err := conn.WriteMessage(messageType, responseData); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}

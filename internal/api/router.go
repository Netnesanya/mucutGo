package api

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

func BroadcastMessage(message string) {
	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Printf("Error broadcasting message to a connection: %v", err)
			// Optionally, handle disconnection
			conn.Close()
			delete(connections, conn)
		}
	}
}

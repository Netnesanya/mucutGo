package api

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/fetch-video-info", VideoInfoHandler).Methods("POST")

	router.HandleFunc("/download-mp3-bulk", Mp3DownloadBulkHandler).Methods("POST")

	router.HandleFunc("/download-mp3", Mp3DownloadHandler).Methods("POST")

	router.HandleFunc("/update-mp3-metadata", Mp3UpdateMetadataHandler).Methods("POST")

	router.HandleFunc("/update-mp3-metadata-bulk", Mp3UpdateMetadataBulkHandler).Methods("POST")

	router.HandleFunc("/pack-siq", CreateSiqHandler).Methods("GET")

	router.HandleFunc("/ws/video-info", videoInfoWebSocket)

	return router
}

func videoInfoWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		// Here, you would read messages from the client, such as a batch of video titles
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// Process the message and send updates back to the client
		// For demonstration, just echo the message back
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust the origin checking for production use
	},
}

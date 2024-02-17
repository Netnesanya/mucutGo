package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/fetch-video-info", VideoInfoHandler).Methods("POST")

	router.HandleFunc("/download-mp3-bulk", Mp3DownloadBulkHandler).Methods("POST")

	router.HandleFunc("/download-mp3", Mp3DownloadHandler).Methods("POST")

	router.HandleFunc("/update-mp3-metadata", Mp3UpdateMetadataHandler).Methods("POST")

	return router
}

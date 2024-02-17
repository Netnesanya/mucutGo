package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/fetch-video-info", VideoInfoHandler).Methods("POST")

	router.HandleFunc("/download-mp3-bulk", Mp3DownloadBulkHandler).Methods("POST")

	return router
}

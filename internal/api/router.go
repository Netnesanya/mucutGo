package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/fetch-video-info", VideoInfoHandler).Methods("POST")

	//router.HandleFunc("/pack-siq", SiqHandler).Methods("GET")

	return router
}

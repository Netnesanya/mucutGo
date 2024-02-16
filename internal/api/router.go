package api

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/video-titles", VideoTitlesHandler).Methods("POST")
	return router
}

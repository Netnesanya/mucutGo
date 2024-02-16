package api

import (
	"encoding/json"
	"io/ioutil"
	"mucutGo/internal/yt"
	"net/http"
	"strings"
)

func VideoTitlesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Assuming the file content is a list of titles, one per line
	titles := strings.Split(string(body), "\n")

	// Prefetch video data
	metadata, err := yt.FetchVideoMetadata(titles)
	if err != nil {
		http.Error(w, "Error prefetching videos", http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(metadata)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func parseTitles(titles string) []string {
	return strings.Split(titles, "\n")
}

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mucutGo/internal/yt"
	"net/http"
	"strings"
)

func VideoInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err) // Log the error
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Splitting titles based on newline
	titles := strings.Split(string(body), "\n")

	// Prefetch video data
	metadata, err := yt.FetchVideoMetadata(titles)
	if err != nil {
		log.Printf("Error prefetching videos for titles %v: %v", titles, err) // Log the error with titles for context
		http.Error(w, fmt.Sprintf("Error prefetching videos: %v", err), http.StatusInternalServerError)
		return
	}

	yt.DownloadAudioFromMetadata(metadata) // Download audio from metadata
	// If no metadata is returned, it might be useful to log this as a separate case
	if len(metadata) == 0 {
		log.Printf("No metadata fetched for titles: %v", titles)                 // This could indicate an issue with yt-dlp or the input titles
		http.Error(w, "No video metadata could be fetched", http.StatusNotFound) // Using 404 might be more descriptive if no data is found
		return
	}

	responseData, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("Error marshaling metadata to JSON for titles %v: %v", titles, err) // Log marshaling error
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

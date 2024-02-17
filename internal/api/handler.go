package api

import (
	"encoding/json"
	"fmt"
	"io"
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

	body, err := io.ReadAll(r.Body)
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

func Mp3DownloadBulkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var metadata []yt.VideoMetadata

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON body into the metadata slice
	if err := json.Unmarshal(body, &metadata); err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		http.Error(w, "Error processing request body", http.StatusBadRequest)
		return
	}

	// Now that you have the metadata, you can pass it to your download function
	err = yt.DownloadAudioFromMetadata(metadata)
	if err != nil {
		log.Printf("Error downloading audio from metadata: %v", err)
		http.Error(w, "Error downloading audio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Audio download initiated successfully"))
}

func Mp3DownloadHandler(w http.ResponseWriter, r *http.Request) {

}

func Mp3UpdateMetadataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var metadata yt.VideoMetadata
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	fmt.Println(string(body)) // Convert the body to a string to print it

	// Unmarshal the request body into the struct
	err = json.Unmarshal(body, &metadata)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	// Process the update (for simplicity, this part is left as a comment)
	// Here, you might use ytsearch with metadata.OriginalUrl or metadata.Title
	// and update the metadata accordingly.
	updatedMetadata, err := yt.FetchMetaDataSingleMp3(metadata)
	if err != nil {
		log.Printf("Error fetching metadata for '%s': %v", metadata.Title, err)
		http.Error(w, "Error fetching metadata", http.StatusInternalServerError)
		return
	}

	// For the sake of this example, let's pretend we've updated the metadata
	// Now, marshal and return the updated metadata
	updatedMetadataJSON, err := json.Marshal(updatedMetadata)
	if err != nil {
		log.Printf("Error marshaling updated metadata: %v", err)
		http.Error(w, "Error marshaling updated metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(updatedMetadataJSON)
}

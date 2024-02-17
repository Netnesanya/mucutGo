package yt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type VideoHeatmap struct {
	EndTime   float32 `json:"end_time"`
	StartTime float32 `json:"start_time"`
	Value     float64 `json:"value"` // Updated type to float64 to handle decimal values
}

type VideoMetadata struct {
	Title          string         `json:"title"`
	Duration       float32        `json:"duration"`        // Assuming duration is in seconds
	DurationString string         `json:"duration_string"` // This might need to be calculated separately if not provided directly
	Heatmap        []VideoHeatmap `json:"heatmap"`
	OriginalUrl    string         `json:"original_url"`
}

// FetchVideoMetadata takes a list of video titles and fetches their metadata using yt-dlp.
func FetchVideoMetadata(titles []string) ([]VideoMetadata, error) {
	var metadataList []VideoMetadata

	for _, title := range titles {
		cmdArgs := []string{
			"--default-search", "ytsearch1:", // Limit to the first search result
			"--dump-json",   // Get the output in JSON format
			"--no-playlist", // Ensure only single video info is returned
			title,           // The search query
		}

		cmd := exec.Command("yt-dlp", cmdArgs...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing yt-dlp for title '%s': %v", title, err)
			continue // Skip this iteration on error
		}

		var videoMeta VideoMetadata
		err = json.Unmarshal(output, &videoMeta)
		if err != nil {
			log.Printf("Error unmarshaling JSON for title '%s': %v", title, err)
			continue // Skip this iteration on error
		}

		metadataList = append(metadataList, videoMeta)
	}

	// Consider error handling if no metadata could be fetched.
	if len(metadataList) == 0 {
		return nil, fmt.Errorf("no metadata could be fetched for the given titles")
	}

	writeMetadataToFile(metadataList, "metadata.json")
	return metadataList, nil
}

func DownloadAudioFromMetadata(metadataList []VideoMetadata) error {
	var errorMessages []string // Collect error messages here

	for _, metadata := range metadataList {
		startTime, endTime := FindHeatmapSpike(metadata.Heatmap, metadata.Duration)
		fileName := fmt.Sprintf("%s.mp3", metadata.Title) // Sanitize to ensure valid filenames
		outputPath := fmt.Sprintf("downloads/%s", fileName)

		log.Printf("Downloading audio segment for '%s'", metadata.Title)
		err := DownloadAudioSegment(metadata.OriginalUrl, startTime, endTime, outputPath)
		if err != nil {
			errorMessage := fmt.Sprintf("Error downloading audio segment for '%s': %v", metadata.Title, err)
			log.Print(errorMessage)
			errorMessages = append(errorMessages, errorMessage)
			continue // Skip this iteration on error
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("errors occurred during downloads: %s", strings.Join(errorMessages, "; "))
	}

	return nil
}

func writeMetadataToFile(metadataList []VideoMetadata, filePath string) error {
	jsonData, err := json.MarshalIndent(metadataList, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	log.Printf("Metadata successfully written to %s\n", filePath)
	return nil
}

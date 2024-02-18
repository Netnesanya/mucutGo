package yt

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type VideoHeatmap struct {
	EndTime   float32 `json:"end_time"`
	StartTime float32 `json:"start_time"`
	Value     float64 `json:"value"` // Updated type to float64 to handle decimal values
}

type UserInput struct {
	From     float32 `json:"from,omitempty"` // Use pointers to make fields optional
	To       float32 `json:"to,omitempty"`
	Duration float32 `json:"duration,omitempty"`
}

type VideoMetadata struct {
	Title          string         `json:"title"`
	Duration       float32        `json:"duration"`        // Assuming duration is in seconds
	DurationString string         `json:"duration_string"` // This might need to be calculated separately if not provided directly
	Heatmap        []VideoHeatmap `json:"heatmap"`
	OriginalUrl    string         `json:"original_url"`
}

type CombinedData struct {
	VideoMetadata VideoMetadata `json:"fetched"`
	UserInput     UserInput     `json:"userInput"`
}

// FetchVideoMetadataFromText takes a list of video titles and fetches their metadata using yt-dlp.
func FetchVideoMetadataFromText(titles []string) ([]VideoMetadata, error) {
	fmt.Println(titles)

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

	return metadataList, nil
}

func DownloadAudioFromMetadata(combinedDataList []CombinedData) error {
	var errorMessages []string // Collect error messages here

	defaultDuration := float32(35) // Default duration in seconds, adjust as needed

	for _, combinedData := range combinedDataList {
		var startTime, endTime float32

		// Check if user input is provided and use it to determine start and end times
		if combinedData.UserInput.From != 0 || combinedData.UserInput.To != 0 {
			startTime = combinedData.UserInput.From
			endTime = combinedData.UserInput.To
			if combinedData.UserInput.Duration != 0 && endTime == 0 {
				// Calculate 'to' using 'from' + 'duration' if 'to' is not provided
				endTime = startTime + combinedData.UserInput.Duration
			}
		} else {
			// Pass userSpecifiedDuration to FindHeatmapSpike, if specified; otherwise, use default duration.
			userSpecifiedDuration := defaultDuration
			if combinedData.UserInput.Duration != 0 {
				userSpecifiedDuration = combinedData.UserInput.Duration
			}
			startTime, endTime = FindHeatmapSpike(combinedData.VideoMetadata.Heatmap, combinedData.VideoMetadata.Duration, &userSpecifiedDuration)
		}

		// Ensure endTime does not exceed video duration
		if endTime == 0 || endTime > combinedData.VideoMetadata.Duration {
			endTime = combinedData.VideoMetadata.Duration
		}

		fileName := fmt.Sprintf("%s.mp3", combinedData.VideoMetadata.Title) // Sanitize to ensure valid filenames
		outputPath := fmt.Sprintf("downloads/%s", fileName)

		log.Printf("Downloading audio segment for '%s' from %.2f to %.2f", combinedData.VideoMetadata.Title, startTime, endTime)
		err := DownloadAudioSegment(combinedData.VideoMetadata.OriginalUrl, startTime, endTime, outputPath)
		if err != nil {
			errorMessage := fmt.Sprintf("Error downloading audio segment for '%s': %v", combinedData.VideoMetadata.Title, err)
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

func FetchMetaDataSingleMp3(metadata VideoMetadata) (VideoMetadata, error) {
	searchQuery := metadata.Title // Default to using the title
	if metadata.OriginalUrl != "" {
		searchQuery = metadata.OriginalUrl // If OriginalUrl is not empty, use it instead
	}

	cmdArgs := []string{
		"--default-search", "ytsearch1:", // Limit to the first search result
		"--dump-json",   // Get the output in JSON format
		"--no-playlist", // Ensure only single video info is returned
		searchQuery,     // Use the selected query
	}

	cmd := exec.Command("yt-dlp", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing yt-dlp for '%s': %v", searchQuery, err)
		return VideoMetadata{}, err // Return an empty VideoMetadata and the error
	}

	var videoMeta VideoMetadata
	err = json.Unmarshal(output, &videoMeta)
	if err != nil {
		log.Printf("Error unmarshaling JSON for '%s': %v", searchQuery, err)
		return VideoMetadata{}, err // Return an empty VideoMetadata and the error
	}

	return videoMeta, nil
}

func FetchMetaDataBulkMp3(metadataArr []VideoMetadata) ([]VideoMetadata, error) {
	var metadataList []VideoMetadata

	for _, metadata := range metadataArr {
		searchQuery := metadata.Title // Default to using the title
		if metadata.OriginalUrl != "" {
			searchQuery = metadata.OriginalUrl // If OriginalUrl is not empty, use it instead
		}

		cmdArgs := []string{
			"--default-search", "ytsearch1:", // Limit to the first search result
			"--dump-json",   // Get the output in JSON format
			"--no-playlist", // Ensure only single video info is returned
			searchQuery,     // The search query
		}

		cmd := exec.Command("yt-dlp", cmdArgs...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing yt-dlp for title '%s': %v", searchQuery, err)
			continue // Skip this iteration on error
		}

		var videoMeta VideoMetadata
		err = json.Unmarshal(output, &videoMeta)
		if err != nil {
			log.Printf("Error unmarshaling JSON for title '%s': %v", searchQuery, err)
			continue // Skip this iteration on error
		}

		metadataList = append(metadataList, videoMeta)
	}

	if len(metadataList) == 0 {
		return nil, fmt.Errorf("no metadata could be fetched for the given titles")
	}

	return metadataList, nil
}

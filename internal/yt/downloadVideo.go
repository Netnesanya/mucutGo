package yt

import (
	"fmt"
	"os/exec"
)

func FindHeatmapSpike(heatmap []VideoHeatmap, duration float32, cutDuration *float32) (startTime, endTime float32) {
	if len(heatmap) == 0 {
		startTime := (duration / 2) - *cutDuration/2
		return startTime, startTime + *cutDuration/2 // Return immediately if heatmap is empty
	}

	var maxSpike VideoHeatmap
	maxSpike.Value = -1.0             // Initialize with a very small value
	ignoreFirstSeconds := float32(20) // Ignore spikes in the first 20 seconds

	for _, point := range heatmap {
		// Skip the initial 20 seconds of the song
		if point.StartTime < ignoreFirstSeconds {
			continue
		}
		// Find the maximum spike after the first 20 seconds
		if point.Value > maxSpike.Value {
			maxSpike = point
		}
	}

	// If no spike found after the first 20 seconds, it might mean all spikes are within the first 20 seconds
	// In such a case, or if maxSpike.Value remains -1, indicating no spike was found, you may need a fallback strategy
	if maxSpike.Value == -1.0 {
		// Fallback strategy: could return the start of the video or another logic
		// For now, let's return the first 30-35 seconds after the 20 seconds mark
		startTime = ignoreFirstSeconds
		endTime = min(startTime+*cutDuration, duration) // Ensure we do not exceed the video duration
		return
	}

	// Start 10 seconds before the spike
	startTimeAdjustment := float32(10)
	startTime = max(maxSpike.StartTime-startTimeAdjustment, 0) // Ensure start time is not negative

	// The end time is determined by adding the cutDuration to the startTime, ensuring it doesn't exceed the video duration
	endTime = min(startTime+*cutDuration, duration) // Ensure we do not exceed the video duration

	return
}

// Helper function to find the minimum of two float32 values
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// Helper function to find the maximum of two float32 values
func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func DownloadAudioSegment(url string, startTime, endTime float32, outputPath string) error {
	if url == "" {
		return fmt.Errorf("URL is empty, cannot download segment")
	}

	fmt.Println(url, startTime, endTime, outputPath)

	cmdArgs := []string{
		url,
		"-x",                    // Extract audio
		"--audio-format", "mp3", // Specify audio format, adjust as needed
		"-o", outputPath, // Specify the output path and filename
		"--external-downloader", "ffmpeg",
		"--external-downloader-args", fmt.Sprintf("ffmpeg_i:-ss %.2f -to %.2f", startTime, endTime),
	}

	cmd := exec.Command("yt-dlp", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download audio segment: %s, error: %v", string(output), err)
	}

	return nil
}

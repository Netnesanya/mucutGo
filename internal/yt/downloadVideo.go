package yt

import (
	"fmt"
	"os/exec"
)

//func DownloadVideo(url string, heatmap []VideoHeatmap) error {
//	startTime, endTime := FindHeatmapSpike(heatmap)
//	if err := DownloadAudioSegment(url, startTime, endTime, "/"); err != nil {
//		return err
//	}
//	return nil
//}

func FindHeatmapSpike(heatmap []VideoHeatmap, duration float32) (startTime, endTime float32) {
	if len(heatmap) == 0 {
		return 0, 0 // Return immediately if heatmap is empty
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
		endTime = min(startTime+30, duration) // Ensure we do not exceed the video duration
		return
	}

	halfCutLength := float32(17.5) // Adjust for a total of 30-35 seconds cut
	// Center the cut around the spike's maximum value, adjusting for video boundaries
	startTime = max(maxSpike.StartTime-halfCutLength, ignoreFirstSeconds) // Respect the initial ignore period
	endTime = min(maxSpike.StartTime+halfCutLength, duration)             // Ensure we do not exceed the video duration

	// Additional adjustment if calculated cut is shorter than desired due to video start constraint
	if endTime-startTime < 35 && endTime < duration {
		additionalTime := min(35-(endTime-startTime), duration-endTime)
		endTime += additionalTime
	}

	return
}

// Helper functions to find the minimum and maximum of two float32 values
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func DownloadAudioSegment(url string, startTime, endTime float32, outputPath string) error {
	// Modify command arguments for audio extraction
	cmdArgs := []string{
		url,
		"-x",                    // Extract audio
		"--audio-format", "mp3", // Specify audio format, adjust as needed
		"-o", outputPath, // Specify the output path and filename
		"--external-downloader", "ffmpeg",
		"--external-downloader-args", fmt.Sprintf("ffmpeg_i:-ss %.2f -to %.2f", startTime, endTime),
	}

	cmd := exec.Command("yt-dlp", cmdArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to download audio segment: %s, error: %v", string(output), err)
	}

	return nil
}

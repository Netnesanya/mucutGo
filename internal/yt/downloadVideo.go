package yt

import (
	"fmt"
	"os/exec"
)

func DownloadVideo(url string, heatmap []VideoHeatmap) error {
	startTime, endTime := FindHeatmapSpike(heatmap)
	if err := DownloadAudioSegment(url, startTime, endTime, "/"); err != nil {
		return err
	}
	return nil
}

func FindHeatmapSpike(heatmap []VideoHeatmap) (startTime, endTime float32) {
	var maxSpike VideoHeatmap
	for _, point := range heatmap {
		if point.Value > maxSpike.Value {
			maxSpike = point
		}
	}

	startTime = maxSpike.StartTime - 5
	if startTime < 0 {
		startTime = 0
	}

	endTime = maxSpike.EndTime + 5
	return
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

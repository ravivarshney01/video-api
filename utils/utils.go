package utils

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func MergeVideos(filePaths []string, outputPath string) error {
	if len(filePaths) == 0 {
		return fmt.Errorf("no input files provided")
	}

	// Create a temporary file list
	tempFile := fmt.Sprintf("merge_temp_file_list_%d.txt", time.Now().Unix())
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile) // Cleanup temporary file after use

	// Write file paths to the temporary file
	for _, filePath := range filePaths {
		_, err := file.WriteString(fmt.Sprintf("file '%s'\n", filePath))
		if err != nil {
			return fmt.Errorf("failed to write to temporary file: %v", err)
		}
	}
	file.Close()

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", tempFile, "-c", "copy", outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("FFmpeg command failed: %v", err)
	}

	return nil
}

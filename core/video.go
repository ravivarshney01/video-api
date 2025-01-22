package core

import (
	"context"
	"fmt"
	"mime/multipart"
	"os/exec"
	"time"
	"video-api/models"
	"video-api/uploadSvc"
	"video-api/utils"
)

func UploadVideo(ctx context.Context, file multipart.File, fileName string) (id int, err error) {
	uploadService := uploadSvc.NewFileSystemUploadService()
	path, err := uploadService.Upload(fileName, file)
	if err != nil {
		return
	}

	id, err = models.AddVideo(ctx, fileName, path)
	if err != nil {
		return
	}
	return
}

func TrimVideo(ctx context.Context, videoId int, start, end int) (err error) {
	video, err := models.GetVideo(ctx, videoId)
	if err != nil {
		return
	}

	outputPath := fmt.Sprintf("uploads/trimmed_%s", video.Filename)
	cmd := exec.Command("ffmpeg", "-ss", fmt.Sprintf("%d", start), "-to", fmt.Sprintf("%d", end), "-i", video.Url, "-c", "copy", outputPath)
	if err := cmd.Run(); err != nil {
		return
	}
	video.Url = outputPath
	video.Filename = "trimmed_" + video.Filename
	err = video.UpdateVideo(ctx)
	if err != nil {
		return
	}
	return
}

func MergeVideos(ctx context.Context, videoIds []int) (id int, err error) {
	videos, err := models.FindVideosByIds(ctx, videoIds)
	if err != nil {
		return
	}

	var videoPaths []string
	for _, video := range videos {
		videoPaths = append(videoPaths, video.Url)
	}

	fileName := fmt.Sprintf("merge_%d.mp4", time.Now().Unix())
	outputPath := fmt.Sprintf("uploads/%s", fileName)
	err = utils.MergeVideos(videoPaths, outputPath)
	if err != nil {
		return
	}

	id, err = models.AddVideo(ctx, fileName, outputPath)
	if err != nil {
		return
	}
	return
}

package models

import (
	"context"
	"video-api/db"
)

type Video struct {
	Id       int    `json:"id"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
}

func AddVideo(ctx context.Context, fileName string, url string) (id int, err error) {
	video := &Video{Filename: fileName, Url: url}
	err = db.Create(ctx, &video)
	return video.Id, err
}

func GetVideo(ctx context.Context, id int) (video Video, err error) {
	err = db.First(ctx, &video, id)
	return video, err
}

func (video *Video) UpdateVideo(ctx context.Context) error {
	return db.Save(ctx, video)
}

func FindVideosByIds(ctx context.Context, ids []int) ([]Video, error) {
	var videos []Video
	err := db.FindByIds(ctx, &videos, ids)
	return videos, err
}

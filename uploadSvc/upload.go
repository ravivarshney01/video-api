package uploadSvc

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"
)

type UploadService interface {
	Upload(fileName string, file multipart.File) (string, error)
}

type FileSystemUploadService struct{}

func (FileSystemUploadService) Upload(fileName string, file multipart.File) (url string, err error) {
	filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), fileName)
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	defer out.Close()
	_, err = out.ReadFrom(file)
	if err != nil {
		return
	}

	url = fmt.Sprintf(filename)
	return
}

func NewFileSystemUploadService() UploadService {
	return &FileSystemUploadService{}
}

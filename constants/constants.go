package constants

import "time"

const (
	MaximumFileToUpload = 10 * 1024 * 1024 //10 MB
	SharedLinkExpiry    = 10 * time.Minute
)

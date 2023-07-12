package global

import "log"

var (
	Log           *log.Logger = nil
	DownloadRetry int         = 5
	TempStorage   string      = "raw"
)

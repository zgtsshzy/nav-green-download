package mfwam

import (
	"context"
	"fmt"
	"nav-green-download/pkg/server"
	"time"
)

func ExecuteScript(startStr, endStr string) error {
	start, err := time.Parse(time.DateTime, startStr)
	if err != nil {
		return fmt.Errorf("start 时间解析失败: %v", err)
	}

	end, err := time.Parse(time.DateTime, endStr)
	if err != nil {
		return fmt.Errorf("end 时间解析失败: %v", err)
	}

	mfwam := server.NewMFWAMDownloader()
	defer mfwam.Stop(context.TODO())

	for !start.After(end) {
		mfwam.DownloadByDate(start)
		start = start.Add(time.Hour * 12)
	}

	return nil
}

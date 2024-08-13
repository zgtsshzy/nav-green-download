package seaice

import (
	"context"
	"fmt"
	"nav-green-download/pkg/server"
	"time"

	"github.com/sirupsen/logrus"
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

	seaIce := server.NewSeaIceDownloader()
	defer seaIce.Stop(context.TODO())

	for !start.After(end) {
		logrus.Infof("开始下载: %v 海冰NC数据文件", start)
		seaIce.DownloadByDate(start)
		start = start.Add(time.Hour * 24)
	}

	return nil
}

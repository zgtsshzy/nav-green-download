package server

import (
	"context"
	"fmt"
	"nav-green-download/pkg/conf"
	"nav-green-download/pkg/global"
	"nav-green-download/pkg/tools"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type GFSDownloader struct {
}

func NewGFSDownloader() *GFSDownloader {
	config := conf.Get()

	if _, err := os.Stat(config.GFSDir); err != nil {
		os.Mkdir(config.GFSDir, 0777)
	}

	return &GFSDownloader{}
}

func (srv *GFSDownloader) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour)
	childCtx := context.WithoutCancel(ctx)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("GFS 文件下载程序停止")
		case <-ticker.C:
			srv.Download(childCtx)
			ticker.Reset(time.Hour)
		}
	}
}

func (srv *GFSDownloader) Download(ctx context.Context) {
	ch := make(chan DownloadInfo, 1)

	go func() {
		if err := srv.getFirstLevelInfo(ch); err != nil {
			logrus.Errorf("获取 GFS 下载文件列表失败: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case info, ok := <-ch:
			if !ok {
				logrus.Infof("下载任务队列关闭")
				return
			}

			if err := tools.DownloadNCFile(info.LocalPath, info.Url); err != nil {
				logrus.Errorf("GFS: %s 下载失败: %v", info.Url, err)
			}
		}
	}
}

func (srv *GFSDownloader) getFirstLevelInfo(ch chan DownloadInfo) error {
	defer close(ch)

	firstLevels, err := tools.GetGFSFirstLevel()
	if err != nil {
		return fmt.Errorf("GetGFSFirstLevel: %v", err)
	}

	for _, first := range firstLevels {
		folder, _ := strings.CutSuffix(first, "/")
		localDir := filepath.Join(conf.Get().GFSDir, folder)
		if _, err := os.Stat(localDir); err != nil {
			if err := os.Mkdir(localDir, 0777); err != nil {
				logrus.Errorf("文件夹: %s 创建失败: %v", localDir, err)
			}
		}

		if err := srv.getSecondLevelInfo(first, ch); err != nil {
			return fmt.Errorf("getSecondLevelInfo: %v", err)
		}
	}

	return nil
}

func (srv *GFSDownloader) getSecondLevelInfo(level string, ch chan DownloadInfo) error {
	secondLevels, err := tools.GetGFSSecondLevel(level)
	if err != nil {
		return fmt.Errorf("GetGFSSecondLevel: %v", err)
	}

	for _, second := range secondLevels {
		folder, _ := strings.CutSuffix(level+second, "/")
		localDir := filepath.Join(conf.Get().GFSDir, folder)
		if _, err := os.Stat(localDir); err != nil {
			if err := os.Mkdir(localDir, 0777); err != nil {
				logrus.Errorf("文件夹: %s 创建失败: %v", localDir, err)
			}
		}

		if err := srv.getThirdLevelInfo(level+second, ch); err != nil {
			return fmt.Errorf("getThirdLevelInfo: %v", err)
		}
	}

	return nil
}

func (srv *GFSDownloader) getThirdLevelInfo(level string, ch chan DownloadInfo) error {
	thirdLevels, err := tools.GetGFSThirdLevel(level)
	if err != nil {
		return fmt.Errorf("GetGFSThirdLevel: %v", err)
	}

	for _, third := range thirdLevels {
		folder, _ := strings.CutSuffix(level+third, "/")
		localDir := filepath.Join(conf.Get().GFSDir, folder)
		if _, err := os.Stat(localDir); err != nil {
			if err := os.Mkdir(localDir, 0777); err != nil {
				logrus.Errorf("文件夹: %s 创建失败: %v", localDir, err)
			}
		}

		if err := srv.getFourthLevelInfo(level+third, ch); err != nil {
			return fmt.Errorf("getFourthLevelInfo: %v", err)
		}
	}

	return nil
}

func (srv *GFSDownloader) getFourthLevelInfo(level string, ch chan DownloadInfo) error {
	fourthLevels, err := tools.GetGFSFourthLevel(level)
	if err != nil {
		return fmt.Errorf("GetGFSFourthLevel: %v", err)
	}

	for _, fourth := range fourthLevels {
		folder, _ := strings.CutSuffix(level+fourth, "/")
		localPath := filepath.Join(conf.Get().GFSDir, folder)
		downloadUrl := fmt.Sprintf("%s%s", global.GFSBaseUrl, level+fourth)

		if _, err := os.Stat(localPath); err != nil {
			ch <- DownloadInfo{
				LocalPath: localPath,
				Url:       downloadUrl,
			}
		}
	}

	return nil
}

func (srv *GFSDownloader) Stop(ctx context.Context) error {

	return nil
}

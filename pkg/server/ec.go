package server

import (
	"context"
	"fmt"
	"nav-green-download/pkg/conf"
	"nav-green-download/pkg/global"
	"nav-green-download/pkg/tools/download"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ECDownloader struct {
}

func NewECDownloader() *ECDownloader {
	config := conf.Get()

	if _, err := os.Stat(config.ECDir); err != nil {
		os.Mkdir(config.ECDir, 0777)
	}

	return &ECDownloader{}
}

func (srv *ECDownloader) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour)
	childCtx := context.WithoutCancel(ctx)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("EC 文件下载程序停止")
		case <-ticker.C:
			srv.Download(childCtx)
			ticker.Reset(time.Hour)
		}
	}
}

type DownloadInfo struct {
	LocalPath string
	Url       string
}

func (srv *ECDownloader) Download(ctx context.Context) {
	ch := make(chan DownloadInfo, 1)

	go func() {
		if err := srv.getFirstLevelInfo(ch); err != nil {
			logrus.Errorf("获取 EC 下载文件列表失败: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case info, ok := <-ch:
			if !ok {
				return
			}

			if err := download.DownloadNCFile(info.LocalPath, info.Url); err != nil {
				logrus.Errorf("EC 文件: %s 下载失败: %v", info.Url, err)
			}
		}
	}
}

func (srv *ECDownloader) getFirstLevelInfo(ch chan DownloadInfo) error {
	defer close(ch)

	firstLevels, err := download.GetECFirstLevel()
	if err != nil {
		return fmt.Errorf("GetECFirstLevel: %v", err)
	}

	for _, first := range firstLevels {
		folder, _ := strings.CutSuffix(first, "/")
		localDir := filepath.Join(conf.Get().ECDir, folder)
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

func (srv *ECDownloader) getSecondLevelInfo(level string, ch chan DownloadInfo) error {
	secondLevels, err := download.GetECSecondLevel(level)
	if err != nil {
		return fmt.Errorf("GetECSecondLevel: %v", err)
	}

	for _, second := range secondLevels {
		folder, _ := strings.CutSuffix(level+second, "/")
		localDir := filepath.Join(conf.Get().ECDir, folder)
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

func (srv *ECDownloader) getThirdLevelInfo(level string, ch chan DownloadInfo) error {
	thirdLevels, err := download.GetECThirdLevel(level)
	if err != nil {
		return fmt.Errorf("GetECThirdLevel: %v", err)
	}

	for _, third := range thirdLevels {
		folder, _ := strings.CutSuffix(level+third, "/")
		localDir := filepath.Join(conf.Get().ECDir, folder)
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

func (srv *ECDownloader) getFourthLevelInfo(level string, ch chan DownloadInfo) error {
	fourthLevels, err := download.GetECFourthLevel(level)
	if err != nil {
		return fmt.Errorf("GetECFourthLevel: %v", err)
	}

	for _, fourth := range fourthLevels {
		folder, _ := strings.CutSuffix(level+fourth, "/")
		localDir := filepath.Join(conf.Get().ECDir, folder)
		if _, err := os.Stat(localDir); err != nil {
			if err := os.Mkdir(localDir, 0777); err != nil {
				logrus.Errorf("文件夹: %s 创建失败: %v", localDir, err)
			}
		}

		if err := srv.getFifthLevelInfo(level+fourth, ch); err != nil {
			return fmt.Errorf("getFifthLevelInfo: %v", err)
		}
	}

	return nil
}

func (srv *ECDownloader) getFifthLevelInfo(level string, ch chan DownloadInfo) error {
	fifthLevels, err := download.GetECFifthLevel(level)
	if err != nil {
		return fmt.Errorf("GetECFifthLevel: %v", err)
	}

	for _, fifth := range fifthLevels {
		folder, _ := strings.CutSuffix(level+fifth, "/")
		localDir := filepath.Join(conf.Get().ECDir, folder)
		if _, err := os.Stat(localDir); err != nil {
			if err := os.Mkdir(localDir, 0777); err != nil {
				logrus.Errorf("文件夹: %s 创建失败: %v", localDir, err)
			}
		}

		if err := srv.getSixthLevelInfo(level+fifth, ch); err != nil {
			return fmt.Errorf("getSixthLevelInfo: %v", err)
		}
	}

	return nil
}

func (srv *ECDownloader) getSixthLevelInfo(level string, ch chan DownloadInfo) error {
	sixthLevels, err := download.GetECSixthFiles(level)
	if err != nil {
		return fmt.Errorf("GetECSixthFiles: %v", err)
	}

	for _, sixth := range sixthLevels {
		folder, _ := strings.CutSuffix(level+sixth, "/")
		localPath := filepath.Join(conf.Get().ECDir, folder)
		downloadUrl := fmt.Sprintf("%s%s", global.ECBaseUrl, level+sixth)

		if _, err := os.Stat(localPath); err != nil {
			ch <- DownloadInfo{
				LocalPath: localPath,
				Url:       downloadUrl,
			}
		}
	}

	return nil
}

func (srv *ECDownloader) Stop(ctx context.Context) error {

	return nil
}

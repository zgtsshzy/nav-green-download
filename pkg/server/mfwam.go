package server

import (
	"context"
	"fmt"
	"nav-green-download/pkg/conf"
	"nav-green-download/pkg/tools"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MFWAMDownloader struct {
}

func NewMFWAMDownloader() *MFWAMDownloader {
	config := conf.Get()

	if _, err := os.Stat(config.MFWAMDir); err != nil {
		os.Mkdir(config.MFWAMDir, 0777)
	}

	return nil
}

func (srv *MFWAMDownloader) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("程序停止")
		case <-ticker.C:
			// 从今天开始下载未来 10 天的 NC 文件
			currentDate, _ := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
			for i := 0; i < 20; i++ {
				srv.DownloadByDate(currentDate)
				currentDate = currentDate.Add(time.Hour * 12)
			}

			// 从今天开始下载过去 10 天的 NC 文件
			currentDate, _ = time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
			for i := 0; i < 20; i++ {
				srv.DownloadByDate(currentDate)
				currentDate = currentDate.Add(time.Hour * -12)
			}
		}
	}
}

func (srv *MFWAMDownloader) Stop(ctx context.Context) error {

	return nil
}

func (srv *MFWAMDownloader) DownloadByDate(date time.Time) error {
	//////////////////////////////////////////////////////////////////////////////
	// 确认文件是否已经下载过，如果和服务器一样则不需要再下载
	localName := getMFWAMLocalNameByDate(date)
	remoteName, err := tools.GetMFWAMNameByDate(date)
	if err != nil {
		return fmt.Errorf("获取官网: %v MFWAM 文件名失败: %v", date, err)
	}

	if localName != "" && localName == remoteName {
		return nil
	}

	//////////////////////////////////////////////////////////////////////////////
	{
		// 删除过期的 MFWAM 文件
		path := getMFWAMLocalPathByDate(date)
		if path != "" {
			os.Remove(path)
		}
	}

	//////////////////////////////////////////////////////////////////////////////
	// 下载最新的 MFWAM 文件到本地
	localPath, err := generateMFWAMLocalPathByDate(date)
	if err != nil {
		return err
	}

	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	url, err := tools.GetMFWAMDownloadUrlByDate(date)
	if err != nil {
		return fmt.Errorf("获取 MFWAM : %v 文件下载链接失败: %v", date, err)
	}

	if err := tools.DownloadNCFile(localPath, url); err != nil {
		return fmt.Errorf("下载 MFWAM : %v 文件失败: %v", date, err)
	}

	return nil
}

func generateMFWAMLocalPathByDate(date time.Time) (string, error) {
	url, err := tools.GetMFWAMDownloadUrlByDate(date)
	if err != nil {
		return "", fmt.Errorf("获取 MFWAM : %v 文件下载链接失败: %v", date, err)
	}

	itemList := strings.Split(url, "/")
	fileName := itemList[len(itemList)-1]

	path := filepath.Join(
		conf.Get().MFWAMDir,
		fmt.Sprintf("%d", date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fileName,
	)

	return path, nil
}

// 根据日期获取 MFWAM 本地文件名称
func getMFWAMLocalNameByDate(date time.Time) string {
	dateStr := date.Format("2006010215")

	dir := filepath.Join(
		conf.Get().MFWAMDir,
		fmt.Sprintf("%d", date.Year()),
		fmt.Sprintf("%02d", date.Month()),
	)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		list := strings.Split(entry.Name(), "_")
		if len(list) != 3 {
			continue
		}

		if list[1] == dateStr {
			return entry.Name()
		}
	}

	return ""
}

// 根据日期获取 MFWAM 本地文件路径
func getMFWAMLocalPathByDate(date time.Time) string {
	dateStr := date.Format("2006010215")

	dir := filepath.Join(
		conf.Get().MFWAMDir,
		fmt.Sprintf("%d", date.Year()),
		fmt.Sprintf("%02d", date.Month()),
	)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		list := strings.Split(entry.Name(), "_")
		if len(list) != 3 {
			continue
		}

		if list[1] == dateStr {
			return filepath.Join(dir, entry.Name())
		}
	}

	return ""
}

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

type SeaIceDownloader struct {
}

func NewSeaIceDownloader() *SeaIceDownloader {
	config := conf.Get()

	if _, err := os.Stat(config.SeaIceDir); err != nil {
		os.Mkdir(config.SeaIceDir, 0777)
	}

	return nil
}

func (srv *SeaIceDownloader) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("程序停止")
		case <-ticker.C:
			// 从今天开始下载未来 10 天的 NC 文件
			currentDate, _ := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
			for i := 0; i < 10; i++ {
				srv.DownloadByDate(currentDate)
				currentDate = currentDate.Add(time.Hour * 24)
			}

			// 从今天开始下载过去 10 天的 NC 文件
			currentDate, _ = time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
			for i := 0; i < 10; i++ {
				srv.DownloadByDate(currentDate)
				currentDate = currentDate.Add(time.Hour * -24)
			}
		}
	}
}

func (srv *SeaIceDownloader) Stop(ctx context.Context) error {

	return nil
}

func (srv *SeaIceDownloader) DownloadByDate(date time.Time) error {
	//////////////////////////////////////////////////////////////////////////////
	// 确认文件是否已经下载过，如果和服务器一样则不需要再下载
	localName := getSeaIceLocalNameByDate(date)
	remoteName, err := tools.GetSeaIceNameByDate(date)
	if err != nil {
		return fmt.Errorf("获取官网: %v SeaIce 文件名失败: %v", date, err)
	}

	if localName != "" && localName == remoteName {
		return nil
	}

	//////////////////////////////////////////////////////////////////////////////
	{
		// 删除过期的 SeaIce 文件
		path := getSeaIceLocalPathByDate(date)
		if path != "" {
			os.Remove(path)
		}
	}

	//////////////////////////////////////////////////////////////////////////////
	// 下载最新的 SeaIce 文件到本地
	localPath, err := generateSeaIceLocalPathByDate(date)
	if err != nil {
		return err
	}

	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	url, err := tools.GetSeaIceDownloadUrlByDate(date)
	if err != nil {
		return fmt.Errorf("获取 SeaIce : %v 文件下载链接失败: %v", date, err)
	}

	if err := tools.DownloadNCFile(localPath, url); err != nil {
		return fmt.Errorf("下载 SeaIce : %v 文件失败: %v", date, err)
	}

	return nil
}

func generateSeaIceLocalPathByDate(date time.Time) (string, error) {
	url, err := tools.GetSeaIceDownloadUrlByDate(date)
	if err != nil {
		return "", fmt.Errorf("获取 SeaIce : %v 文件下载链接失败: %v", date, err)
	}

	itemList := strings.Split(url, "/")
	fileName := itemList[len(itemList)-1]

	path := filepath.Join(
		conf.Get().SeaIceDir,
		fmt.Sprintf("%d", date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fileName,
	)

	return path, nil
}

// 根据日期获取 SeaIce 本地文件名称
func getSeaIceLocalNameByDate(date time.Time) string {
	dir := filepath.Join(
		conf.Get().SeaIceDir,
		fmt.Sprintf("%d", date.Year()),
		fmt.Sprintf("%02d", date.Month()),
	)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		list := strings.Split(entry.Name(), "_")
		if len(list) != 7 {
			continue
		}

		dateStr := strings.Split(list[1], "-")[0]
		if dateStr == date.Format("20060102") {
			return entry.Name()
		}
	}

	return ""
}

// 根据日期获取 SeaIce 本地文件路径
func getSeaIceLocalPathByDate(date time.Time) string {
	dir := filepath.Join(
		conf.Get().SeaIceDir,
		fmt.Sprintf("%d", date.Year()),
		fmt.Sprintf("%02d", date.Month()),
	)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		// glo12_rg_1d-m_20240801-20240801_2D_hcst_R20240814.nc
		list := strings.Split(entry.Name(), "_")
		if len(list) != 7 {
			continue
		}

		dateStr := strings.Split(list[1], "-")[0]
		if dateStr == date.Format("20060102") {
			return filepath.Join(dir, entry.Name())
		}
	}

	return ""
}

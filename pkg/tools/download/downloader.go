package download

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

var httpClient *resty.Client

func init() {
	httpClient = resty.New()
	httpClient.SetRetryCount(5)
	httpClient.SetTimeout(time.Minute * 10)
}

// 根据 URL 下载 NC 文件，保存到本地路径下
func DownloadNCFile(localPath, url string) error {
	os.Remove(localPath)

	_, err := httpClient.R().SetOutput(localPath).Get(url)
	if err != nil {
		return fmt.Errorf("下载 nc 文件: %s 失败: %v", localPath, err)
	}

	return nil
}

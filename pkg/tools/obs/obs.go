package obs

import (
	"fmt"
	"time"
)

type OBSInfo struct {
	DateTime  time.Time // 时间
	Path      string    // 本地的文件路径
	UploadKey string    // OBS 上传的KEY
}

// 上传文件到 OBS
func UploadFile2OBS(info OBSInfo) error {
	if info.Path == "" {
		return fmt.Errorf("本地文件路径为空")
	}

	return nil
}

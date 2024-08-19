package download

import (
	"fmt"
	"nav-green-download/pkg/global"
	"strings"

	"github.com/antchfx/htmlquery"
)

// 获取 GFS 官网第一级页面信息
// gfs.20240807/           07-Aug-2024 20:41
// gfs.20240808/           08-Aug-2024 20:41
// gfs.20240809/           09-Aug-2024 20:41
// gfs.20240810/           10-Aug-2024 20:41
// gfs.20240811/           11-Aug-2024 20:41
// gfs.20240812/           12-Aug-2024 20:41
// gfs.20240813/           13-Aug-2024 20:41
// gfs.20240814/           14-Aug-2024 20:41
// gfs.20240815/           15-Aug-2024 20:41
// gfs.20240816/           16-Aug-2024 14:41
func GetGFSFirstLevel() ([]string, error) {
	htmlNode, err := htmlquery.LoadURL(global.GFSBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("ec 第一级页面获取失败: %v", err)
	}

	var folders []string
	valueList := htmlquery.Find(htmlNode, `/html/body/pre/a`)
	for _, value := range valueList {
		name := htmlquery.InnerText(value)
		if strings.Contains(name, "gfs") {
			folders = append(folders, name)
		}
	}

	return folders, nil
}

// 获取 GFS 官网第二级页面信息
// 00/                          14-Aug-2024 03:33
// 06/                          14-Aug-2024 09:29
// 12/                          14-Aug-2024 15:30
// 18/                          14-Aug-2024 21:35
func GetGFSSecondLevel(level string) ([]string, error) {
	htmlNode, err := htmlquery.LoadURL(global.GFSBaseUrl + level)
	if err != nil {
		return nil, fmt.Errorf("gfs 第二级页面获取失败: %v", err)
	}

	var folders []string
	valueList := htmlquery.Find(htmlNode, `/html/body/pre/a`)
	for _, value := range valueList {
		name := htmlquery.InnerText(value)
		if name == "00/" {
			folders = append(folders, name)
		}
	}

	return folders, nil
}

// 获取 GFS 官网第三级页面信息
// atmos/                                    14-Aug-2024 05:15
// wave/                                     14-Aug-2024 04:47
func GetGFSThirdLevel(level string) ([]string, error) {
	htmlNode, err := htmlquery.LoadURL(global.GFSBaseUrl + level)
	if err != nil {
		return nil, fmt.Errorf("gfs 第三级页面获取失败: %v", err)
	}

	var folders []string
	valueList := htmlquery.Find(htmlNode, `/html/body/pre/a`)
	for _, value := range valueList {
		name := htmlquery.InnerText(value)
		if name == "atmos/" {
			folders = append(folders, name)
		}
	}

	return folders, nil
}

// 获取 GFS 官网第四级页面信息
// gfs.t00z.pgrb2b.0p25.f237                    14-Aug-2024 04:37
// gfs.t00z.pgrb2b.0p25.f240                    14-Aug-2024 04:37
// gfs.t00z.pgrb2b.0p25.f243                    14-Aug-2024 04:39
// gfs.t00z.pgrb2b.0p25.f246                    14-Aug-2024 04:39
// gfs.t00z.pgrb2b.0p25.f249                    14-Aug-2024 04:40
// gfs.t00z.pgrb2b.0p25.f252                    14-Aug-2024 04:41
// gfs.t00z.pgrb2b.0p25.f255                    14-Aug-2024 04:42
func GetGFSFourthLevel(level string) ([]string, error) {
	htmlNode, err := htmlquery.LoadURL(global.GFSBaseUrl + level)
	if err != nil {
		return nil, fmt.Errorf("gfs 第四级页面获取失败: %v", err)
	}

	var folders []string
	valueList := htmlquery.Find(htmlNode, `/html/body/pre/a`)
	for _, value := range valueList {
		name := htmlquery.InnerText(value)
		if strings.Contains(name, ".idx") || strings.HasSuffix(name, ".anl") {
			continue
		}

		if strings.Contains(name, ".pgrb2b.0p25.") {
			folders = append(folders, htmlquery.InnerText(value))
		}
	}

	return folders, nil
}

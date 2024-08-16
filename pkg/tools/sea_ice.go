package tools

import (
	"encoding/xml"
	"fmt"
	"nav-green-download/pkg/global"
	"strings"
	"time"
)

type SeaIceResult struct {
	XMLName  xml.Name      `xml:"ListBucketResult"`
	Contents []SMOCContent `xml:"Contents"`
}

type SeaIceContent struct {
	Key string `xml:"Key"`
}

func GetSeaIceDownloadUrlByDate(date time.Time) (string, error) {
	prefix := fmt.Sprintf("native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_0.083deg_P1D-m_202406/%d/%02d/", date.Year(), date.Month())

	resp, err := httpClient.R().
		SetQueryParams(map[string]string{
			"delimiter": "/",
			"list-type": "2",
			"prefix":    prefix,
		}).
		Get("https://mdl-native-14.s3.waw3-1.cloudferro.com/")

	if err != nil {
		return "", fmt.Errorf("获取 SeaIce 官网: %v 下载列表失败: %v", date, err)
	}

	var result SeaIceResult
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("unmarshal SeaIce XML列表失败: %v", err)
	}

	for _, content := range result.Contents {
		// https://s3.waw3-1.cloudferro.com/mdl-native-14/native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_0.083deg_P1D-m_202406/2024/08/glo12_rg_1d-m_20240801-20240801_2D_hcst_R20240814.nc
		//                                                native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_0.083deg_P1D-m_202406/2024/08/glo12_rg_1d-m_20240801-20240801_2D_hcst_R20240814.nc
		// glo12_rg_1d-m_20240801-20240801_2D_hcst_R20240814.nc
		nameList := strings.Split(content.Key, "/")
		name := nameList[len(nameList)-1]

		dateStr := strings.Split(name, "_")[3]
		dateStr = strings.Split(dateStr, "-")[0]

		if date.Format("20060102") == dateStr {
			return global.SeaIceBaseUrl + content.Key, nil
		}
	}

	return "", fmt.Errorf("SeaIce: %v 没有找到下载链接", date)
}

func GetSeaIceNameByDate(date time.Time) (string, error) {
	prefix := fmt.Sprintf("native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_0.083deg_P1D-m_202406/%d/%02d/", date.Year(), date.Month())

	resp, err := httpClient.R().
		SetQueryParams(map[string]string{
			"delimiter": "/",
			"list-type": "2",
			"prefix":    prefix,
		}).
		Get("https://mdl-native-14.s3.waw3-1.cloudferro.com/")

	if err != nil {
		return "", fmt.Errorf("获取 SeaIce 官网: %v 下载列表失败: %v", date, err)
	}

	var result SeaIceResult
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("unmarshal SeaIce XML列表失败: %v", err)
	}

	for _, content := range result.Contents {
		// https://s3.waw3-1.cloudferro.com/mdl-native-14/native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_0.083deg_P1D-m_202406/2024/08/glo12_rg_1d-m_20240801-20240801_2D_hcst_R20240814.nc
		//                                                native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_0.083deg_P1D-m_202406/2024/08/glo12_rg_1d-m_20240801-20240801_2D_hcst_R20240814.nc
		// glo12_rg_1d-m_20240801-20240801_2D_hcst_R20240814.nc
		nameList := strings.Split(content.Key, "/")
		name := nameList[len(nameList)-1]

		dateStr := strings.Split(name, "_")[3]
		dateStr = strings.Split(dateStr, "-")[0]

		if date.Format("20060102") == dateStr {
			return name, nil
		}
	}

	return "", fmt.Errorf("SeaIce: %v 没有找到下载链接", date)
}

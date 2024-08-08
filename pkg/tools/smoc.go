package tools

import (
	"encoding/xml"
	"fmt"
	"nav-green-download/pkg/global"
	"strings"
	"time"
)

type SMOCResult struct {
	XMLName  xml.Name      `xml:"ListBucketResult"`
	Contents []SMOCContent `xml:"Contents"`
}

type SMOCContent struct {
	Key string `xml:"Key"`
}

func GetSMOCDownloadUrlByDate(date time.Time) (string, error) {
	prefix := fmt.Sprintf("native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211/%d/%02d/", date.Year(), date.Month())

	resp, err := httpClient.R().
		SetQueryParams(map[string]string{
			"delimiter": "/",
			"list-type": "2",
			"prefix":    prefix,
		}).
		Get("https://mdl-native-14.s3.waw3-1.cloudferro.com/")

	if err != nil {
		return "", fmt.Errorf("获取 SMOC 官网: %v 下载列表失败: %v", date, err)
	}

	var result SMOCResult
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("unmarshal SMOC XML列表失败: %v", err)
	}

	for _, content := range result.Contents {
		// https://s3.waw3-1.cloudferro.com/mdl-native-14/native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211/2024/08/SMOC_20240808_R20240730.nc
		//                                                native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211/2024/06/SMOC_20240602_R20240612.nc
		nameList := strings.Split(content.Key, "/")
		name := nameList[len(nameList)-1]

		dateStr := strings.Split(name, "_")[1]

		if date.Format("20060102") == dateStr {
			return global.SMOCBaseUrl + content.Key, nil
		}
	}

	return "", fmt.Errorf("SMOC: %v 没有找到下载链接", date)
}

func GetSMOCNameByDate(date time.Time) (string, error) {
	prefix := fmt.Sprintf("native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211/%d/%02d/", date.Year(), date.Month())

	resp, err := httpClient.R().
		SetQueryParams(map[string]string{
			"delimiter": "/",
			"list-type": "2",
			"prefix":    prefix,
		}).
		Get("https://mdl-native-14.s3.waw3-1.cloudferro.com/")

	if err != nil {
		return "", fmt.Errorf("获取 SMOC 官网: %v 下载列表失败: %v", date, err)
	}

	var result SMOCResult
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("unmarshal SMOC XML列表失败: %v", err)
	}

	for _, content := range result.Contents {
		// https://s3.waw3-1.cloudferro.com/mdl-native-14/native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211/2024/08/SMOC_20240808_R20240730.nc
		//                                                native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211/2024/06/SMOC_20240602_R20240612.nc
		nameList := strings.Split(content.Key, "/")
		name := nameList[len(nameList)-1]

		dateStr := strings.Split(name, "_")[1]

		if date.Format("20060102") == dateStr {
			return name, nil
		}
	}

	return "", fmt.Errorf("SMOC: %v 没有找到下载链接", date)
}

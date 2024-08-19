package download

import (
	"encoding/xml"
	"fmt"
	"nav-green-download/pkg/global"
	"strings"
	"time"
)

type MFWAMResult struct {
	XMLName  xml.Name       `xml:"ListBucketResult"`
	Contents []MFWAMContent `xml:"Contents"`
}

type MFWAMContent struct {
	Key string `xml:"Key"`
}

// https://data.marine.copernicus.eu/product/GLOBAL_ANALYSISFORECAST_WAV_001_027/files?subdataset=cmems_mod_glo_wav_anfc_0.083deg_PT3H-i_202311&path=GLOBAL_ANALYSISFORECAST_WAV_001_027%2Fcmems_mod_glo_wav_anfc_0.083deg_PT3H-i_202311%2F
func GetMFWAMDownloadUrlByDate(date time.Time) (string, error) {
	prefix := fmt.Sprintf("native/GLOBAL_ANALYSISFORECAST_WAV_001_027/cmems_mod_glo_wav_anfc_0.083deg_PT3H-i_202311/%d/%02d/", date.Year(), date.Month())

	resp, err := httpClient.R().
		SetQueryParams(map[string]string{
			"delimiter": "/",
			"list-type": "2",
			"prefix":    prefix,
		}).
		Get("https://mdl-native-14.s3.waw3-1.cloudferro.com/")

	if err != nil {
		return "", fmt.Errorf("获取 MFWAM 官网: %v 下载列表失败: %v", date, err)
	}

	var result MFWAMResult
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("unmarshal MFWAM XML列表失败: %v", err)
	}

	for _, content := range result.Contents {
		nameList := strings.Split(content.Key, "/")
		name := nameList[len(nameList)-1]

		dateStr := strings.Split(name, "_")[1]

		if date.Format("2006010215") == dateStr {
			return global.MFWAMBaseUrl + content.Key, nil
		}
	}

	return "", fmt.Errorf("MFWAM: %v 没有找到下载链接", date)
}

func GetMFWAMNameByDate(date time.Time) (string, error) {
	prefix := fmt.Sprintf("native/GLOBAL_ANALYSISFORECAST_WAV_001_027/cmems_mod_glo_wav_anfc_0.083deg_PT3H-i_202311/%d/%02d/", date.Year(), date.Month())

	resp, err := httpClient.R().
		SetQueryParams(map[string]string{
			"delimiter": "/",
			"list-type": "2",
			"prefix":    prefix,
		}).
		Get("https://mdl-native-14.s3.waw3-1.cloudferro.com/")

	if err != nil {
		return "", fmt.Errorf("获取 MFWAM 官网: %v 下载列表失败: %v", date, err)
	}

	var result MFWAMResult
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("unmarshal MFWAM XML列表失败: %v", err)
	}

	for _, content := range result.Contents {
		// native/GLOBAL_ANALYSISFORECAST_WAV_001_027/cmems_mod_glo_wav_anfc_0.083deg_PT3H-i_202311/2024/08/mfwamglocep_2024080112_R20240802_12H.nc
		nameList := strings.Split(content.Key, "/")
		name := nameList[len(nameList)-1]

		dateStr := strings.Split(name, "_")[1]

		if date.Format("2006010215") == dateStr {
			return name, nil
		}
	}

	return "", fmt.Errorf("MFWAM: %v 没有找到下载链接", date)
}

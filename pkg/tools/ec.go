package tools

import (
	"fmt"
	"nav-green-download/pkg/global"

	"github.com/antchfx/htmlquery"
)

// https://www.ecmwf.int/en/forecasts/datasets/open-data
// EC 气象数据下载官网
// https://data.ecmwf.int/forecasts/
// https://data.ecmwf.int/forecasts/20240815/00z/ifs/0p25/enfo/20240815000000-0h-enfo-ef.grib2

// 获取 EC 官网第一级页面信息
// 20240812/      12-08-2024 07:55
// 20240813/      13-08-2024 07:55
// 20240814/      14-08-2024 07:55
// 20240815/      15-08-2024 07:55
func GetECFirstLevel() ([]string, error) {
	htmlNode, err := htmlquery.LoadURL(global.ECBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("ec 第一级页面获取失败: %v", err)
	}

	var folders []string
	valueList := htmlquery.Find(htmlNode, `//*[@id="outerTable"]/tbody/tr[4]/td/pre[2]/a`)
	for _, value := range valueList {
		folders = append(folders, htmlquery.InnerText(value))
	}

	return folders, nil
}

// 获取 EC 官网第二级页面信息
// 00z/           12-08-2024 07:55
// 06z/           12-08-2024 13:12
// 12z/           12-08-2024 19:55
// 18z/    	      13-08-2024 01:12
func GetECSecondLevel(level string) ([]string, error) {
	htmlNode, err := htmlquery.LoadURL(global.ECBaseUrl + level)
	if err != nil {
		return nil, fmt.Errorf("ec 第二级页面获取失败: %v", err)
	}

	var folders []string
	valueList := htmlquery.Find(htmlNode, `//*[@id="outerTable"]/tbody/tr[4]/td/pre[2]/a`)
	for _, value := range valueList {
		folders = append(folders, htmlquery.InnerText(value))
	}

	return folders, nil
}

// 获取 EC 官网第三级页面信息
// aifs/          12-08-2024 07:55
// ifs/           12-08-2024 08:40
func GetECThirdLevel(level string) ([]string, error) {
	return []string{"ifs/"}, nil
}

// 获取 EC 官网第四级页面信息
// 0p25/          12-08-2024 08:40
// 0p4-beta/      12-08-2024 08:40
func GetECFourthLevel(level string) ([]string, error) {
	return []string{"0p25/"}, nil
}

// 获取 EC 官网第五级页面信息
// enfo/          12-08-2024 08:40
// oper/          12-08-2024 07:55
// waef/          12-08-2024 08:40
// wave/          12-08-2024 07:55
func GetECFifthLevel(level string) ([]string, error) {
	htmlNode, err := htmlquery.LoadURL(global.ECBaseUrl + level)
	if err != nil {
		return nil, fmt.Errorf("ec 第五级页面获取失败: %v", err)
	}

	var folders []string
	valueList := htmlquery.Find(htmlNode, `//*[@id="outerTable"]/tbody/tr[4]/td/pre[2]/a`)
	for _, value := range valueList {
		folders = append(folders, htmlquery.InnerText(value))
	}

	return folders, nil
}

// 获取 EC 官网第六级页面信息
// 20240812000000-0h-enfo-ef.grib2       12-08-2024 08:40    4233252678      87258125
// 20240812000000-0h-enfo-ef.index       12-08-2024 08:40       1297700      87258131
// 20240812000000-102h-enfo-ef.grib2     12-08-2024 08:40    4484172942      87259929
// 20240812000000-102h-enfo-ef.index     12-08-2024 08:40       1310196      87259935
// 20240812000000-105h-enfo-ef.grib2     12-08-2024 08:40    4489709581      87259440
// 20240812000000-105h-enfo-ef.index     12-08-2024 08:40       1310300      87259443
func GetECSixthFiles(level string) ([]string, error) {
	htmlNode, err := htmlquery.LoadURL(global.ECBaseUrl + level)
	if err != nil {
		return nil, fmt.Errorf("ec 第六级页面获取失败: %v", err)
	}

	var folders []string
	valueList := htmlquery.Find(htmlNode, `//*[@id="outerTable"]/tbody/tr[4]/td/pre[2]/a`)
	for _, value := range valueList {
		folders = append(folders, htmlquery.InnerText(value))
	}

	return folders, nil
}

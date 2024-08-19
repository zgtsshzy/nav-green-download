package download

import "testing"

func TestDownloadNCFile(t *testing.T) {
	localPath := "./SMOC_20240808_R20240730.nc"
	url := "https://s3.waw3-1.cloudferro.com/mdl-native-14/native/GLOBAL_ANALYSISFORECAST_PHY_001_024/cmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211/2024/08/SMOC_20240808_R20240730.nc"

	if err := DownloadNCFile(localPath, url); err != nil {
		t.Errorf("DownloadNCFile: %v", err)
		return
	}
}

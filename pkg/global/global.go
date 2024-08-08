package global

const (
	// 日志配置
	DefaultLevel      = "info"
	DefaultMaxLogSize = 10
	DefaultMaxLogAge  = 7
	DefaultMaxBackups = 5

	// https://data.marine.copernicus.eu/product/GLOBAL_ANALYSISFORECAST_PHY_001_024/files?path=GLOBAL_ANALYSISFORECAST_PHY_001_024%2Fcmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211%2F2024%2F07%2F&subdataset=cmems_mod_glo_phy_anfc_merged-uv_PT1H-i_202211
	SMOCBaseUrl = "https://s3.waw3-1.cloudferro.com/mdl-native-14/"

	// https://data.marine.copernicus.eu/product/GLOBAL_ANALYSISFORECAST_WAV_001_027/files?subdataset=cmems_mod_glo_wav_anfc_0.083deg_PT3H-i_202311&path=GLOBAL_ANALYSISFORECAST_WAV_001_027%2Fcmems_mod_glo_wav_anfc_0.083deg_PT3H-i_202311%2F
	// https://s3.waw3-1.cloudferro.com/mdl-native-14/native/GLOBAL_ANALYSISFORECAST_WAV_001_027/cmems_mod_glo_wav_anfc_0.083deg_PT3H-i_202311/2024/07/mfwamglocep_2024070100_R20240702_00H.nc?x-cop-client=MyOcean&x-cop-usage=FileBrowser
	MFWAMBaseUrl = "https://s3.waw3-1.cloudferro.com/mdl-native-14/"

	ECBaseUrl = "https://data.ecmwf.int/forecasts/"
)

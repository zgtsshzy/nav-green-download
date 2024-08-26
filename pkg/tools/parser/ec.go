package parser

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"os/exec"
	"path"
	"reflect"

	"github.com/batchatco/go-native-netcdf/netcdf"
	"github.com/batchatco/go-native-netcdf/netcdf/api"
	"github.com/sirupsen/logrus"

	"nav-green-download/pkg/tools/isoline"
)

// var _ Parser = (*ECParser)(nil)

var sfc_dims = []string{"time", "number", "latitude", "longitude"}

var pl_dims = []string{"time", "number", "level", "latitude", "longitude"}

var sfc = []string{"t2m", "u10", "v10", "tp", "d2m", "msl", "ssrd", "lsm", "cape", "v100", "u100", "sp", "swvl1", "swvl2", "ssr", "st", "stl2", "skt", "swvl3", "swvl4", "strd", "stl3", "asn", "tcwv", "stl4", "str", "ro", "ttr"}

var pl = []string{"d", "gh", "q", "r", "t", "u", "v", "w", "vo"}

type ECParser struct {
	GribPath  string
	SFCPath   string
	PLPath    string
	TempDir   string
	OutputDir string
	Del       chan int
	Err       error
}

type ECParserBuilder struct {
	ECParser
}

func NewECParser() (*ECParser, error) {
	return new(ECParser), nil
}

// builder模式构造

func NewECParserBuilder() *ECParserBuilder {
	return new(ECParserBuilder)
}

func (eb *ECParserBuilder) WithGribPath(s string) *ECParserBuilder {
	eb.GribPath = s
	return eb
}

func (eb *ECParserBuilder) WithOutputPath(s string) *ECParserBuilder {
	eb.OutputDir = s
	return eb
}

func (eb *ECParserBuilder) WithTempPath(s string) *ECParserBuilder {
	eb.TempDir = s
	return eb
}

func (eb *ECParserBuilder) WithSFCPath(s string) *ECParserBuilder {
	eb.SFCPath = s
	return eb
}

func (eb *ECParserBuilder) WithPLPath(s string) *ECParserBuilder {
	eb.PLPath = s
	return eb
}

func (eb *ECParserBuilder) Build() *ECParser {
	if len(eb.GribPath) == 0 {
		eb.GribPath = "/data1/test/wd/20240821000000-0h-enfo-ef.grib2"
	}
	if len(eb.TempDir) == 0 {
		eb.TempDir = path.Dir(eb.GribPath)
	}
	if len(eb.OutputDir) == 0 {
		filename := path.Base(eb.GribPath)
		filesuffix := path.Ext(filename)
		fileprefix := eb.GribPath[0 : len(eb.GribPath)-len(filesuffix)]
		eb.OutputDir = fileprefix + "_img/"
	}
	return &eb.ECParser
}

// 实现功能

func (parser *ECParser) SetGribPath(s string) *ECParser {
	// 改变输入文件位置
	parser.GribPath = s
	return parser
}

func (parser *ECParser) RunCmd(name string, arg ...string) *ECParser {
	if parser.Err == nil {
		cmd := exec.Command(name, arg...)
		parser.Err = cmd.Run()
	}
	return parser
}

func (parser *ECParser) ParseData() *ECParser {
	// 解析文件,将文件分为两部分,然后转换成nc文件

	if parser.Err == nil {
		if parser.Del != nil {
			// 重新解析文件将删除原本的临时文件
			close(parser.Del)
			parser.Del = nil
		}

		if _, err := os.Stat(parser.TempDir); err != nil {
			os.Mkdir(parser.TempDir, 0777)
		}

		filename := path.Base(parser.GribPath)
		filesuffix := path.Ext(filename)
		fileprefix := filename[0 : len(filename)-len(filesuffix)]

		if filesuffix != ".grib2" {
			parser.Err = fmt.Errorf("not grib file")
			return parser
		}

		prefix := path.Join(parser.TempDir, fileprefix)

		sfcpath := prefix + "-SFC" + filesuffix
		plpath := prefix + "-PL" + filesuffix
		parser.SFCPath = prefix + "-SFC" + ".nc"
		parser.PLPath = prefix + "-PL" + ".nc"

		// 将原.grib2文件分割为两部分
		parser = parser.RunCmd("grib_copy", "-w", "levtype=sfc", parser.GribPath, sfcpath).RunCmd("grib_copy", "-w", "levtype=pl", parser.GribPath, plpath)
		// 将分割得到的两个.grib2文件转换为两个nc文件
		parser = parser.RunCmd("grib_to_netcdf", "-o", parser.SFCPath, sfcpath).RunCmd("grib_to_netcdf", "-o", parser.PLPath, plpath)

		// 删除分割得到的临时girb2文件
		go func() {
			os.Remove(sfcpath)
			os.Remove(plpath)
		}()

		parser.Del = make(chan int)
		// channel关闭时删除临时文件(nc文件)
		go func() {
			<-parser.Del
			err := os.Remove(parser.SFCPath)
			fmt.Println(err)
			logrus.Warn(err)
			parser.SFCPath = ""
			err = os.Remove(parser.PLPath)
			logrus.Warn(err)
			parser.PLPath = ""
		}()
	}
	return parser
}

func (parser *ECParser) DelTempFile() *ECParser {
	// 手动删除临时文件

	if parser.Del != nil {
		close(parser.Del)
		parser.Del = nil
	}
	return parser
}

func (parser *ECParser) getData4dint16(ds *api.Group, variable string) [][][][]int16 {
	// 从文件中读取数据
	vals, _ := (*ds).GetVariable(variable)
	data, has := vals.Values.([][][][]int16)
	if !has {
		logrus.Warn(fmt.Sprintf("%s has no data", variable))
		return make([][][][]int16, 0)
	}
	return data
}

func (parser *ECParser) getData5dint16(ds *api.Group, variable string) [][][][][]int16 {
	// 从文件中读取数据
	vals, _ := (*ds).GetVariable(variable)
	data, has := vals.Values.([][][][][]int16)
	if !has {
		logrus.Warn(fmt.Sprintf("%s has no data", variable))
		return make([][][][][]int16, 0)
	}
	return data
}

func (parser *ECParser) getData1dint32(ds *api.Group, variable string) []int32 {
	// 从文件中读取数据
	vals, _ := (*ds).GetVariable(variable)
	data, has := vals.Values.([]int32)
	if !has {
		logrus.Warn(fmt.Sprintf("%s has no data", variable))
		return make([]int32, 0)
	}
	return data
}

func (parser *ECParser) getData1dfloat32(ds *api.Group, variable string) []float32 {
	// 从文件中读取数据
	vals, _ := (*ds).GetVariable(variable)
	data, has := vals.Values.([]float32)
	if !has {
		logrus.Warn(fmt.Sprintf("%s has no data", variable))
		return make([]float32, 0)
	}
	return data
}

func (parser *ECParser) drawGray(data [][]int16, outputPath string) {
	img := image.NewGray(image.Rect(0, 0, len(data[0]), len(data)))
	var mi int16 = -32768
	var ma int16 = 32767
	for _, row := range data {
		for _, v := range row {
			mi = min(mi, v)
			ma = max(ma, v)
		}
	}
	bot := float64(mi)
	ran := float64(ma - mi)
	for y, row := range data {
		for x, v := range row {
			gray := math.Floor((float64(v) - bot) / ran * 255)
			img.Set(x, y, color.Gray{Y: uint8(gray)})
		}
	}
	os.Remove(outputPath)
	file, _ := os.Create(outputPath)
	jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
	file.Close()
}

func (parser *ECParser) DrawGrayscale() *ECParser {
	// 做灰度图

	if parser.Err == nil {

		outDir := path.Join(parser.OutputDir, "gray")

		if _, err := os.Stat(outDir); err != nil {
			os.MkdirAll(outDir, 0777)
		}

		if len(parser.SFCPath) == 0 && len(parser.PLPath) == 0 {
			parser.Err = fmt.Errorf("no nc file created")
			return parser
		}

		var ds1, ds2 api.Group

		// sfc部分
		ds1, parser.Err = netcdf.Open(parser.SFCPath)
		if parser.Err != nil {
			return parser
		}
		defer ds1.Close()

		DimSFC := make([][]int32, 0)
		LatLonSFC := make([][]float32, 0)
		for _, dim := range sfc_dims {
			if dim == "latitude" || dim == "longitude" {
				LatLonSFC = append(LatLonSFC, parser.getData1dfloat32(&ds1, dim))
			} else {
				DimSFC = append(DimSFC, parser.getData1dint32(&ds1, dim))
			}
		}

		if len(LatLonSFC[0]) != 721 || len(LatLonSFC[1]) != 1440 {
			logrus.Warn(fmt.Sprintf("unnoraml lat lon: %d, %d", len(LatLonSFC[0]), len(LatLonSFC[1])))
		}

		for _, key := range sfc {
			data4 := parser.getData4dint16(&ds1, key)
			for i, data3 := range data4 {
				for j, data2 := range data3 {
					suffix := fmt.Sprintf("%s_%d_%d_gray.png", key, DimSFC[0][i], DimSFC[1][j])
					outputPath := path.Join(outDir, suffix)
					parser.drawGray(data2, outputPath)
				}
			}
		}

		// pl部分
		ds2, parser.Err = netcdf.Open(parser.PLPath)
		if parser.Err != nil {
			return parser
		}
		defer ds2.Close()

		DimPL := make([][]int32, 0)
		LatLonPL := make([][]float32, 0)
		for _, dim := range pl_dims {
			if dim == "latitude" || dim == "longitude" {
				LatLonPL = append(LatLonPL, parser.getData1dfloat32(&ds2, dim))
			} else {
				DimPL = append(DimPL, parser.getData1dint32(&ds2, dim))
			}
		}

		if len(LatLonPL[0]) != 721 || len(LatLonPL[1]) != 1440 {
			logrus.Warn(fmt.Sprintf("unnoraml lat lon: %d, %d", len(LatLonPL[0]), len(LatLonPL[1])))
		}

		for _, key := range pl {
			data5 := parser.getData5dint16(&ds2, key)
			for i, data4 := range data5 {
				for j, data3 := range data4 {
					for k, data2 := range data3 {
						suffix := fmt.Sprintf("%s_%d_%d_%d_gray.png", key, DimPL[0][i], DimPL[1][j], DimPL[2][k])
						outputPath := path.Join(outDir, suffix)
						parser.drawGray(data2, outputPath)
					}
				}
			}
		}

	}
	return parser
}

func (parser *ECParser) drawIsoline(lats []float32, lons []float32, data [][]int16, config ContourLineConf, outputPathImg string, outputPathJson string) {
	// 画图+生成json
	var start, end, step float32

	start = float32(data[0][0])
	end = float32(data[0][0])

	vals := make([][]float32, 0)
	for i, row := range data {
		vals = append(vals, make([]float32, 0))
		for _, val := range row {
			v := float32(val)
			vals[i] = append(vals[i], v)
			start = min(start, v)
			end = max(end, v)
		}
	}

	if reflect.DeepEqual(config, ContourLineConf{}) {
		config.Num = 10
	}

	if config.MinValue < config.MaxValue && config.Step > 0 {
		start = config.MinValue
		end = config.MaxValue
		step = config.Step
	} else if config.Num > 0 {
		step = (end - start - 1) / float32(config.Num)
	} else {
		logrus.Error(fmt.Errorf("incorrect format of ContourLineConf"))
	}

	LatLonList, thresholds := isoline.FetchRawData(&lats, &lons, &vals, start, end, step)

	// 生成geojson
	err := isoline.SaveGeoJson(isoline.GenGeoJsonFromLatLonList(LatLonList, thresholds), outputPathJson)
	if err != nil {
		logrus.Warn(err)
	}

	// 作图
	x, y, _ := isoline.GenXYList(LatLonList, thresholds)
	err = isoline.MakePlot(x, y, outputPathImg)
	if err != nil {
		logrus.Warn(err)
	}
}

func (parser *ECParser) DrawContourLine(configs map[string]ContourLineConf) *ECParser {
	// 做等值线

	if parser.Err == nil {
		outDir := path.Join(parser.OutputDir, "contour")

		if _, err := os.Stat(outDir); err != nil {
			os.MkdirAll(outDir, 0777)
		}

		if len(parser.SFCPath) == 0 && len(parser.PLPath) == 0 {
			parser.Err = fmt.Errorf("no nc file created")
			return parser
		}

		var ds1, ds2 api.Group

		// sfc部分
		ds1, parser.Err = netcdf.Open(parser.SFCPath)
		if parser.Err != nil {
			return parser
		}
		defer ds1.Close()

		DimSFC := make([][]int32, 0)
		LatLonSFC := make([][]float32, 0)
		for _, dim := range sfc_dims {
			if dim == "latitude" || dim == "longitude" {
				LatLonSFC = append(LatLonSFC, parser.getData1dfloat32(&ds1, dim))
			} else {
				DimSFC = append(DimSFC, parser.getData1dint32(&ds1, dim))
			}
		}

		if len(LatLonSFC[0]) != 721 || len(LatLonSFC[1]) != 1440 {
			logrus.Warn(fmt.Sprintf("unnoraml lat lon: %d, %d", len(LatLonSFC[0]), len(LatLonSFC[1])))
		}

		for _, key := range sfc {
			data4 := parser.getData4dint16(&ds1, key)
			config := configs[key]
			for i, data3 := range data4 {
				for j, data2 := range data3 {
					suffix_img := fmt.Sprintf("%s_%d_%d_contour.png", key, DimSFC[0][i], DimSFC[1][j])
					suffix_json := fmt.Sprintf("%s_%d_%d_contour.json", key, DimSFC[0][i], DimSFC[1][j])
					outputPath_img := path.Join(outDir, suffix_img)
					outputPath_json := path.Join(outDir, suffix_json)

					parser.drawIsoline(LatLonSFC[0], LatLonSFC[1], data2, config, outputPath_img, outputPath_json)
				}
			}
		}

		// pl部分
		ds2, parser.Err = netcdf.Open(parser.PLPath)
		if parser.Err != nil {
			return parser
		}
		defer ds2.Close()

		DimPL := make([][]int32, 0)
		LatLonPL := make([][]float32, 0)
		for _, dim := range pl_dims {
			if dim == "latitude" || dim == "longitude" {
				LatLonPL = append(LatLonPL, parser.getData1dfloat32(&ds2, dim))
			} else {
				DimPL = append(DimPL, parser.getData1dint32(&ds2, dim))
			}
		}

		if len(LatLonPL[0]) != 721 || len(LatLonPL[1]) != 1440 {
			logrus.Warn(fmt.Sprintf("unnoraml lat lon: %d, %d", len(LatLonPL[0]), len(LatLonPL[1])))
		}

		for _, key := range pl {
			data5 := parser.getData5dint16(&ds2, key)
			config := configs[key]

			for i, data4 := range data5 {
				for j, data3 := range data4 {
					for k, data2 := range data3 {

						suffix_img := fmt.Sprintf("%s_%d_%d_%d_contour.png", key, DimPL[0][i], DimPL[1][j], DimPL[2][k])
						suffix_json := fmt.Sprintf("%s_%d_%d_%d_contour.json", key, DimPL[0][i], DimPL[1][j], DimPL[2][k])
						outputPath_img := path.Join(outDir, suffix_img)
						outputPath_json := path.Join(outDir, suffix_json)

						parser.drawIsoline(LatLonSFC[0], LatLonSFC[1], data2, config, outputPath_img, outputPath_json)
					}
				}
			}
		}

	}
	return parser
}

func (parser *ECParser) Save2DB(tableName string) *ECParser {
	if parser.Err == nil {

	}
	return parser
}

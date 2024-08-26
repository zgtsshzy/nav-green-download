package parser

import (
	"fmt"
	"testing"
)

func TestECGray(t *testing.T) {
	parser := NewECParserBuilder().Build()
	defer parser.DelTempFile()
	parser = parser.ParseData().DrawGrayscale()
	if parser.Err != nil {
		fmt.Println(parser.Err)
	}
}

func TestDrawGrayFromNC(t *testing.T) {
	parser := NewECParserBuilder().WithSFCPath("/data1/test/wd/nc_files/20240821000000-0h-enfo-ef-SFC.nc").WithPLPath("/data1/test/wd/nc_files/20240821000000-0h-enfo-ef-PL.nc").Build()
	parser = parser.DrawGrayscale()
	if parser.Err != nil {
		fmt.Println(parser.Err)
	}
}

func TestDrawContourFromNC(t *testing.T) {
	parser := NewECParserBuilder().WithSFCPath("/data1/test/wd/nc_files/20240821000000-0h-enfo-ef-SFC.nc").WithPLPath("/data1/test/wd/nc_files/20240821000000-0h-enfo-ef-PL.nc").WithOutputPath("/data1/test/wd/nc_files/20240821000000-0h-enfo-ef_img/").Build()
	parser = parser.DrawContourLine(make(map[string]ContourLineConf))
	if parser.Err != nil {
		fmt.Println(parser.Err)
	}
}

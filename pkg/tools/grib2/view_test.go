package grib2

import "testing"

func TestLengthDataViewSFC(t *testing.T) {
	// LengthDataViewSFC("/data1/test/20240821000000-0h-enfo-ef-sfc.nc")
	LengthDataViewSFC("/data1/test/wd/20240821000000-0h-enfo-ef-SFC.nc")
}

func TestLengthDataViewPL(t *testing.T) {
	LengthDataViewPL("/data1/test/20240821000000-0h-enfo-ef-pl.nc")
}

func TestLengthTypeViewSFC(t *testing.T) {
	TypeViewSFC("/data1/test/20240821000000-0h-enfo-ef-sfc.nc")
}

func TestLengthTypeViewPL(t *testing.T) {
	TypeViewPL("/data1/test/20240821000000-0h-enfo-ef-pl.nc")
}

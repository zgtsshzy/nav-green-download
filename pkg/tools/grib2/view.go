package grib2

import (
	"fmt"
	"reflect"

	// "github.com/fhs/go-netcdf/netcdf"
	"github.com/batchatco/go-native-netcdf/netcdf"
)

var sfc_dims = []string{"time", "number", "latitude", "longitude"}

var pl_dims = []string{"time", "number", "level", "latitude", "longitude"}

var sfc = []string{"t2m", "u10", "v10", "tp", "d2m", "msl", "ssrd", "lsm", "cape", "v100", "u100", "sp", "swvl1", "swvl2", "ssr", "st", "stl2", "skt", "swvl3", "swvl4", "strd", "stl3", "asn", "tcwv", "stl4", "str", "ro", "ttr"}

var pl = []string{"d", "gh", "q", "r", "t", "u", "v", "w", "vo"}

func LengthDataView(path string, dims []string, keys []string) {
	// 读取nc文件的数据

	// ds, err := netcdf.OpenFile(path, netcdf.NOWRITE)
	ds, err := netcdf.Open(path)
	if err != nil {
		panic(err)
	}
	defer ds.Close()

	fmt.Println(ds.ListVariables())

	read1dfloat32 := func(variable string) {
		vals, _ := ds.GetVariable(variable)
		if vals == nil {
			panic(variable + "data not found")
		}

		data, has := vals.Values.([]float32)
		if !has {
			panic("type err")
		}

		fmt.Println(variable + ":")
		fmt.Println(len(data), data[0], data[len(data)-1])
	}

	read1dint32 := func(variable string) {
		vals, _ := ds.GetVariable(variable)
		if vals == nil {
			panic(variable + "data not found")
		}

		data, has := vals.Values.([]int32)
		if !has {
			panic("type err")
		}

		fmt.Println(variable + ":")
		// fmt.Println(len(data), data[0], data[len(data)-1])
		fmt.Println(len(data), data)

	}

	read5dint16 := func(variable string) {
		vals, _ := ds.GetVariable(variable)
		if vals == nil {
			panic(variable + "data not found")
		}

		data, has := vals.Values.([][][][][]int16)
		if !has {
			panic("type err")
		}
		fmt.Println(variable + ":")
		fmt.Println(len(data), len(data[0]), len(data[0][0]), len(data[0][0][0]), len(data[0][0][0][0]))
		fmt.Println(data[0][0][0][0][0])

	}

	read4dint16 := func(variable string) {
		vals, _ := ds.GetVariable(variable)
		if vals == nil {
			panic(variable + "data not found")
		}

		data, has := vals.Values.([][][][]int16)
		if !has {
			panic("type err")
		}
		fmt.Println(variable + ":")
		fmt.Println(len(data), len(data[0]), len(data[0][0]), len(data[0][0][0]))
		fmt.Println(data[0][0][0][0])

	}
	for _, dim := range dims {
		if dim == "latitude" || dim == "longitude" {
			read1dfloat32(dim)
		} else {
			read1dint32(dim)
		}
	}

	for _, key := range keys {
		if len(dims) == 5 {
			read5dint16(key)
		} else {
			read4dint16(key)
		}
	}

}

func TypeView(path string, dims []string, keys []string) {
	ds, err := netcdf.Open(path)
	if err != nil {
		panic(err)
	}
	defer ds.Close()

	fmt.Println(ds.ListVariables())

	readtest := func(variable string) {
		vals, _ := ds.GetVariable(variable)
		if vals == nil {
			panic(variable + "data not found")
		}

		t := reflect.TypeOf(vals.Values)

		fmt.Println(variable, t)
	}

	for _, dim := range dims {
		readtest(dim)
	}

	for _, v := range keys {
		readtest(v)
	}
}

func TypeViewSFC(path string) {
	TypeView(path, sfc_dims, sfc)
}

func TypeViewPL(path string) {
	TypeView(path, pl_dims, pl)
}

func LengthDataViewSFC(path string) {
	LengthDataView(path, sfc_dims, sfc)
}

func LengthDataViewPL(path string) {
	LengthDataView(path, pl_dims, pl)
}

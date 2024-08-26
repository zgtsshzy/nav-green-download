package isoline

import (
	"bytes"
	"math/rand"
	"os"
	"sync"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func MakePlot(x [][][]float64, y [][][]float64, path string) error {
	if path == "" {
		path = "./test.png"
	}

	n := len(x)
	color_list := make([]drawing.Color, n)
	for i := 0; i < n; i++ {
		color_list[i] = drawing.Color{
			R: uint8(rand.Int() % 256),
			G: uint8(rand.Int() % 256),
			B: uint8(rand.Int() % 256),
			A: 255,
		}
	}
	// fmt.Println(color_list)

	graph := chart.Chart{Series: []chart.Series{}}

	for i := range x {
		for j := range x[i] {
			graph.Series = append(graph.Series, chart.ContinuousSeries{
				XValues: x[i][j],
				YValues: y[i][j],
				Style: chart.Style{
					Hidden:      false,
					StrokeColor: color_list[i],
				},
			})
		}
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return err
	}
	// fmt.Println(buffer.Bytes())
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	file.Write(buffer.Bytes())
	defer file.Close()

	return nil
}

func FetchRawData(lats *[]float32, lons *[]float32, vals *[][]float32, start float32, end float32, step float32) (*[][][]LatLon, []float64) {

	thresholds := make([]float64, 0)
	threshold := start

	for {
		thresholds = append(thresholds, float64(threshold))
		threshold += step
		if threshold > end {
			break
		}
	}

	var wg sync.WaitGroup

	LatLonList := make([][][]LatLon, 0)
	for _, threshold := range thresholds {
		wg.Add(1)
		go func() {
			_, results := IsolineWithThreshold(lats, lons, vals, float32(threshold))
			LatLonList = append(LatLonList, *results)
			wg.Done()
		}()
	}
	wg.Wait()

	return &LatLonList, thresholds
}

func GenXYList(LatLonList *[][][]LatLon, thresholds []float64) ([][][]float64, [][][]float64, []float64) {

	x := make([][][]float64, 0)
	y := make([][][]float64, 0)

	for _, LatLons := range *LatLonList {

		temp_level_x := make([][]float64, 0)
		temp_level_y := make([][]float64, 0)
		for i := range LatLons {
			temp_x := make([]float64, 0)
			temp_y := make([]float64, 0)
			for j := range (LatLons)[i] {
				temp_x = append(temp_x, float64((LatLons)[i][j].Lon))
				temp_y = append(temp_y, float64((LatLons)[i][j].Lat))
			}
			temp_level_x = append(temp_level_x, temp_x)
			temp_level_y = append(temp_level_y, temp_y)
		}
		x = append(x, temp_level_x)
		y = append(y, temp_level_y)
	}

	return x, y, thresholds
}

func FetchData(lats *[]float32, lons *[]float32, vals *[][]float32, start float32, end float32, step float32) ([][][]float64, [][][]float64, []float64) {

	return GenXYList(FetchRawData(lats, lons, vals, start, end, step))
}

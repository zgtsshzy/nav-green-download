package isoline

import (
	"os"

	geojson "github.com/paulmach/go.geojson"
)

func GenGeoJson(lats *[]float32, lons *[]float32, vals *[][]float32, start float32, end float32, step float32) *geojson.FeatureCollection {
	return GenGeoJsonFromLatLonList(FetchRawData(lats, lons, vals, start, end, step))
}

func GenGeoJsonFromLatLonList(LatLonList *[][][]LatLon, thresholds []float64) *geojson.FeatureCollection {

	fc := geojson.NewFeatureCollection()

	for k, LatLons := range *LatLonList {
		for i := range LatLons {
			g := geojson.NewLineStringGeometry(make([][]float64, 0))
			for j := range (LatLons)[i] {
				g.LineString = append(g.LineString, []float64{float64((LatLons)[i][j].Lon), float64((LatLons)[i][j].Lat)})
			}
			f := geojson.NewFeature(g)
			f.SetProperty("value", thresholds[k])
			fc.AddFeature(f)
		}
	}
	return fc
}

func SaveGeoJson(fc *geojson.FeatureCollection, out_path string) error {

	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(out_path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	file.Write(rawJSON)
	defer file.Close()

	return nil
}

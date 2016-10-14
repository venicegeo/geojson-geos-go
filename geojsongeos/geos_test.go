/*
Copyright 2016, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package geojsongeos

import (
	"io/ioutil"
	"testing"

	"github.com/paulsmith/gogeos/geos"
	"github.com/venicegeo/geojson-go/geojson"
)

var inputGeojsonFiles2 = [...]string{
	"test/point.geojson",
	"test/point2.geojson",
	"test/point3.geojson",
	"test/linestring.geojson",
	"test/polygon.geojson",
	"test/multipoint.geojson",
	"test/multilinestring.geojson",
	"test/multipolygon.geojson",
	"test/geometrycollection.geojson",
	"test/featureCollection.geojson"}

var inputWKTFiles = [...]string{
	"test/point.wkt",
	"test/linestring.wkt",
	"test/polygon.wkt",
	"test/multipoint.wkt",
	"test/multilinestring.wkt",
	"test/multipolygon.wkt",
	"test/geometrycollection.wkt"}

func TestMain(t *testing.T) {
	var (
		bytes    []byte
		err      error
		gj       interface{}
		geometry *geos.Geometry
	)
	for inx, fileName := range inputGeojsonFiles2 {
		if gj, err = geojson.ParseFile(fileName); err == nil {
			GeosFromGeoJSON(gj)
			t.Log(inx)
		} else {
			t.Error(err)
		}
	}
	for inx2, fileName2 := range inputWKTFiles {
		if bytes, err = ioutil.ReadFile(fileName2); err == nil {
			if geometry, err = geos.FromWKT(string(bytes)); err == nil {
				_, err = GeoJSONFromGeos(geometry)
				_, err = PointCloud(geometry)
				t.Log(inx2)
			}
		}
		if err != nil {
			t.Error(err)
		}
	}
}

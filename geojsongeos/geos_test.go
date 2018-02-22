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
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	pmgeojson "github.com/paulmach/go.geojson"
	"github.com/paulsmith/gogeos/geos"
	"github.com/venicegeo/geojson-go/geojson"
)

var inputGeometries = [...]string{
	"test/point",
	"test/linestring",
	"test/polygon",
	"test/multipoint",
	"test/multilinestring",
	"test/multipolygon",
	"test/geometrycollection",
}

var inputGeojsonFiles2 = [...]string{
	"test/point.geojson",
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
	"test/multipoint.wkt",
	"test/polygon.wkt",
	"test/multilinestring.wkt",
	"test/multipolygon.wkt",
	"test/geometryCollection.wkt"}

func TestPMGeoJSON(t *testing.T) {
	var (
		bytes    []byte
		err      error
		geometry *geos.Geometry
	)
	// Test all geojsonfiles on GeosFromGeoJSON
	for inx, fileName := range inputGeometries {
		var pmg pmgeojson.Geometry
		if bytes, err = ioutil.ReadFile(fileName + ".geojson"); err == nil {
			if err = json.Unmarshal(bytes, &pmg); err == nil {
				geometry, err = GeosFromGeoJSON(&pmg)
				t.Log(inx)
				if wktbytes, err := ioutil.ReadFile(fileName + ".wkt"); err == nil {
					if wktgeom, err := geos.FromWKT(string(wktbytes)); err == nil {
						if eq, err := wktgeom.Equals(geometry); err == nil {
							if !eq {
								geometrywkt, _ := geometry.ToWKT()
								if string(wktbytes) != geometrywkt {
									t.Errorf("expected (%s), got (%s)", wktbytes, geometrywkt)
								}
							}
						} else {
							t.Error(err)
						}
					} else {
						t.Error(err)
					}
				} else {
					//t.Error(err)
				}
			} else {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}
	}

}
func TestMain(t *testing.T) {
	var (
		bytes    []byte
		err      error
		gj       interface{}
		geometry *geos.Geometry
	)
	// Test all geojsonfiles on GeosFromGeoJSON
	for inx, fileName := range inputGeojsonFiles2 {
		if gj, err = geojson.ParseFile(fileName); err == nil {
			geometry, err = GeosFromGeoJSON(gj)
			t.Log(inx)
		} else {
			t.Error(err)
		}
		//Test round trip on GeoJSONFromGeos
		if err == nil && strings.Compare(fileName, "test/featureCollection.geojson") != 0 {
			_, err = GeoJSONFromGeos(geometry)
			t.Log(GeoJSONFromGeos(geometry))
		}
		if err != nil {
			t.Error(err)
			t.Log(fileName)
		}
	}
	// Test all wkt files on GeoJSONFromGeos aswell as test getPointSlice
	for inx2, fileName2 := range inputWKTFiles {
		if bytes, err = ioutil.ReadFile(fileName2); err == nil {
			if geometry, err = geos.FromWKT(string(bytes)); err == nil {
				gj, err = GeoJSONFromGeos(geometry)
				_, err = getPointSlice(geometry)
				//Test round trip on GeosFromGeoJSON
				if err == nil && strings.Compare(fileName2, "test/multipoint.wkt") != 0 {
					_, err = GeosFromGeoJSON(gj)
				}
				t.Log(geometry)
				t.Log(inx2)
				t.Log(fileName2)
				//test PointCloud
				_, err = PointCloud(geometry)
			}
		}
		if err != nil {
			t.Error(err)
		}
	}
}

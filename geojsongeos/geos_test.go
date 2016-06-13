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

func TestMain(t *testing.T) {
	var (
		bytes    []byte
		err      error
		gj       interface{}
		geometry *geos.Geometry
	)
	filename := "test/test.geojson"
	if gj, err = geojson.ParseFile(filename); err == nil {
		GeosFromGeoJSON(gj)
	} else {
		t.Error(err)
	}
	filename = "test/test.wkt"
	if bytes, err = ioutil.ReadFile(filename); err == nil {
		if geometry, err = geos.FromWKT(string(bytes)); err == nil {
			_, err = GeoJSONFromGeos(geometry)
		}
	}
	if err != nil {
		t.Error(err)
	}
}

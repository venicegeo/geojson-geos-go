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

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/paulsmith/gogeos/geos"
	"github.com/venicegeo/geojson-go/geojson"
)

func parseCoord(input []float64) geos.Coord {
	return geos.NewCoord(input[0], input[1])
}
func parseCoordArray(input [][]float64) []geos.Coord {
	var result []geos.Coord
	for inx := 0; inx < len(input); inx++ {
		result = append(result, parseCoord(input[inx]))
	}
	return result
}

func parseGeoJSON(jsonBytes []byte) {
	gj, err := geojson.Parse(jsonBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	var geometries []interface{}
	switch gjObject := gj.(type) {
	case geojson.FeatureCollection:
		for inx := 0; inx < len(gjObject.Features); inx++ {
			geometries = append(geometries, gjObject.Features[inx].Geometry)
		}
	case geojson.Feature:
		geometries = append(geometries, gjObject.Geometry)
	default:
		geometries = append(geometries, gjObject)
	}
	// // FeatureCollection
	// var features []geojson.Feature
	// var geometries []interface{}
	// err := json.Unmarshal(jsonBytes, &fc)
	// if err != nil {
	// 	// log.Fatal("Unable to unmarshal input: %v", err.Error())
	// }
	// if fc.Type == "FeatureCollection" {
	// 	features = fc.Features
	// 	for inx := 0; inx < len(features); inx++ {
	// 		geometries = append(geometries, features[inx].Geometry)
	// 	}
	// } else {
	// 	var feature geojson.Feature
	// 	err = json.Unmarshal(jsonBytes, &feature)
	// 	if feature.Type == "Feature" {
	// 		// geometries = append(geometries, feature.GetGeometry())
	// 	} else {
	// 		var geometry geojson.Geometry
	// 		err = json.Unmarshal(jsonBytes, &geometry)
	// 		geometries = append(geometries, geometry)
	// 	}
	// }

	for inx := 0; inx < len(geometries); inx++ {
		var geometry *geos.Geometry

		switch gt := geometries[inx].(type) {
		case geojson.Point:
			geometry, err = geos.NewPoint(parseCoord(gt.Coordinates))
		case geojson.LineString:
			geometry, err = geos.NewLineString(parseCoordArray(gt.Coordinates)...)
		case geojson.Polygon:
			var coords []geos.Coord
			var coordsArray [][]geos.Coord
			for jnx := 0; jnx < len(gt.Coordinates); jnx++ {
				coords = parseCoordArray(gt.Coordinates[jnx])
				coordsArray = append(coordsArray, coords)
			}
			geometry, err = geos.NewPolygon(coordsArray[0], coordsArray[1:]...)
		case geojson.MultiPoint:
			var points []*geos.Geometry
			var point *geos.Geometry
			for jnx := 0; jnx < len(gt.Coordinates); jnx++ {
				point, err = geos.NewPoint(parseCoord(gt.Coordinates[jnx]))
				points = append(points, point)
			}
			geometry, err = geos.NewCollection(geos.MULTIPOINT, points...)
		case geojson.MultiLineString:
			var lineStrings []*geos.Geometry
			var lineString *geos.Geometry
			for jnx := 0; jnx < len(gt.Coordinates); jnx++ {
				lineString, err = geos.NewLineString(parseCoordArray(gt.Coordinates[jnx])...)
				lineStrings = append(lineStrings, lineString)
			}
			geometry, err = geos.NewCollection(geos.MULTILINESTRING, lineStrings...)

		// switch gt := geometries[inx].(type) {
		// case map[string]interface{}:
		// 	gtype := gt["type"]
		// 	coordinates := gt["coordinates"]
		// 	switch gtype {
		// 	case "Point":
		// 		switch arrayType := coordinates.(type) {
		// 		case []interface{}:
		// 			geometry, err = geos.NewPoint(parseCoord(arrayType))
		// 		default:
		// 			log.Printf("unexpected type %T\n", arrayType)
		// 		}
		// 	case "LineString":
		// switch arrayType := coordinates.(type) {
		// case []interface{}:
		// 	coords := parseCoordArray(arrayType)
		// 	geometry, err = geos.NewLineString(coords...)
		// 		default:
		// 			log.Printf("unexpected type %T\n", arrayType)
		// 		}
		// 	case "MultiPoint":
		// 		switch arrayType := coordinates.(type) {
		// 		case []interface{}:
		// 			var points []*geos.Geometry
		// 			var point *geos.Geometry
		// 			for inx := 0; inx < len(arrayType); inx++ {
		// 				point, err = geos.NewPoint(parseCoord(arrayType[inx].([]interface{})))
		// 				points = append(points, point)
		// 			}
		// 			geometry, err = geos.NewCollection(geos.MULTIPOINT, points...)
		// 		default:
		// 			log.Printf("unexpected type %T\n", arrayType)
		// 		}
		// 	case "Polygon":
		// 		switch arrayType := coordinates.(type) {
		// 		case []interface{}:
		// var coords []geos.Coord
		// var coordsArray [][]geos.Coord
		// 			for inx := 0; inx < len(arrayType); inx++ {
		// 	coords = parseCoordArray(arrayType[inx].([]interface{}))
		// 	coordsArray = append(coordsArray, coords)
		// }
		// geometry, err = geos.NewPolygon(coordsArray[0], coordsArray[1:]...)
		// 		default:
		// 			log.Printf("unexpected type %T\n", arrayType)
		// 		}
		// 	case "MultiLine":
		// 		switch arrayType := coordinates.(type) {
		// 		case []interface{}:
		// 			var lineStrings []*geos.Geometry
		// 			var lineString *geos.Geometry
		// 			for inx := 0; inx < len(arrayType); inx++ {
		// 				coords := parseCoordArray(arrayType[inx].([]interface{}))
		// 				lineString, err = geos.NewLineString(coords...)
		// 				lineStrings = append(lineStrings, lineString)
		// 			}
		// 			geometry, err = geos.NewCollection(geos.MULTILINESTRING, lineStrings...)
		// 		default:
		// 			log.Printf("unexpected type %T\n", arrayType)
		// 		}
		case geojson.GeometryCollection:
			log.Printf("Unimplemented GeometryCollection")
		case geojson.MultiPolygon:
			log.Printf("Unimplemented MultiPolygon")

		default:
			log.Printf("unexpected type %T\n", gt)
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		if geometry != nil {
			fmt.Print(geometry.String())
		}
	}
}
func main() {
	var args = os.Args[1:]
	var filename string
	if len(args) > 0 {
		filename = args[0]
	} else {
		filename = "test/sample.geojson"
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	parseGeoJSON(file)
}

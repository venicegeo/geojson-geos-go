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
	"errors"
	"fmt"

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

// GeosFromGeoJSON takes a GeoJSON object and returns a GEOS geometry
func GeosFromGeoJSON(input interface{}) (*geos.Geometry, error) {
	var (
		geometry *geos.Geometry
		err      error
	)

	switch gt := input.(type) {
	case *geojson.Point:
		geometry, err = geos.NewPoint(parseCoord(gt.Coordinates))
	case *geojson.LineString:
		geometry, err = geos.NewLineString(parseCoordArray(gt.Coordinates)...)
	case *geojson.Polygon:
		var coords []geos.Coord
		var coordsArray [][]geos.Coord
		for inx := 0; inx < len(gt.Coordinates); inx++ {
			coords = parseCoordArray(gt.Coordinates[inx])
			coordsArray = append(coordsArray, coords)
		}
		geometry, err = geos.NewPolygon(coordsArray[0], coordsArray[1:]...)
	case *geojson.MultiPoint:
		var points []*geos.Geometry
		var point *geos.Geometry
		for jnx := 0; jnx < len(gt.Coordinates); jnx++ {
			point, err = geos.NewPoint(parseCoord(gt.Coordinates[jnx]))
			points = append(points, point)
		}
		geometry, err = geos.NewCollection(geos.MULTIPOINT, points...)
	case *geojson.MultiLineString:
		var lineStrings []*geos.Geometry
		var lineString *geos.Geometry
		for jnx := 0; jnx < len(gt.Coordinates); jnx++ {
			lineString, err = geos.NewLineString(parseCoordArray(gt.Coordinates[jnx])...)
			lineStrings = append(lineStrings, lineString)
		}
		geometry, err = geos.NewCollection(geos.MULTILINESTRING, lineStrings...)

	case *geojson.GeometryCollection:
		err = errors.New("Unimplemented GeometryCollection in GeosFromGeoJSON")
	case *geojson.MultiPolygon:
		var (
			coords      []geos.Coord
			coordsArray [][]geos.Coord
			polygons    []*geos.Geometry
			polygon     *geos.Geometry
		)
		for _, polygonCoords := range gt.Coordinates {
			for _, ringCoords := range polygonCoords {
				coords = parseCoordArray(ringCoords)
				coordsArray = append(coordsArray, coords)
			}
			if polygon, err = geos.NewPolygon(coordsArray[0], coordsArray[1:]...); err != nil {
				return nil, err
			}
			polygons = append(polygons, polygon)

		}
		if geometry, err = geos.NewCollection(geos.MULTIPOLYGON, polygons...); err != nil {
			return nil, err
		}
	case *geojson.Feature:
		return GeosFromGeoJSON(gt.Geometry)
	case *geojson.FeatureCollection:
		var (
			geometries []*geos.Geometry
			gc         *geos.Geometry
		)
		for _, feature := range gt.Features {
			if geometry, err = GeosFromGeoJSON(feature); err != nil {
				return nil, err
			}
			geometries = append(geometries, geometry)
		}
		if gc, err = geos.NewCollection(geos.GEOMETRYCOLLECTION, geometries...); err != nil {
			return nil, err
		}
		if geometry, err = gc.Buffer(0); err != nil {
			return nil, err
		}

		return geometry, nil
	default:
		err = fmt.Errorf("Unexpected type in GeosFromGeoJSON: %T\n", gt)
	}
	return geometry, err
}

// GeoJSONFromGeos takes a GEOS geometry and returns a GeoJSON object
func GeoJSONFromGeos(input *geos.Geometry) (interface{}, error) {
	var (
		result interface{}
		err    error
		gType  geos.GeometryType
		coords []geos.Coord
	)
	gType, err = input.Type()
	if err == nil {
		switch gType {
		case geos.POINT:
			var xval, yval float64
			if xval, err = input.X(); err != nil {
				return nil, err
			}
			if yval, err = input.Y(); err != nil {
				return nil, err
			}
			result = geojson.NewPoint([]float64{xval, yval})
		case geos.LINESTRING:
			if coords, err = input.Coords(); err != nil {
				return nil, err
			}
			result = geojson.NewLineString(arrayFromCoords(coords))
		case geos.POLYGON:
			var (
				coordinates [][][]float64
				ring        *geos.Geometry
				rings       []*geos.Geometry
			)
			if ring, err = input.Shell(); err != nil {
				return nil, err
			}
			if coords, err = ring.Coords(); err != nil {
				return nil, err
			}
			coordinates = append(coordinates, arrayFromCoords(coords))
			if rings, err = input.Holes(); err != nil {
				return nil, err
			}
			for _, ring = range rings {

				if coords, err = ring.Coords(); err != nil {
					return nil, err
				}
				coordinates = append(coordinates, arrayFromCoords(coords))
			}
			result = geojson.NewPolygon(coordinates)

		case geos.MULTIPOLYGON:
			var (
				count       int
				coordinates [][][][]float64
				polygon     *geos.Geometry
				polygonIfc  interface{}
				gjPolygon   *geojson.Polygon
				ok          bool
			)
			if count, err = input.NGeometry(); err != nil {
				return nil, err
			}
			for inx := 0; inx < count; inx++ {
				if polygon, err = input.Geometry(inx); err != nil {
					return nil, err
				}
				polygonIfc, err = GeoJSONFromGeos(polygon)
				if gjPolygon, ok = polygonIfc.(*geojson.Polygon); !ok {
					return nil, fmt.Errorf("Expected Polygon, received %T", polygonIfc)
				}
				coordinates = append(coordinates, gjPolygon.Coordinates)
			}
			result = geojson.NewMultiPolygon(coordinates)
		default:
			err = fmt.Errorf("Unimplemented %v", gType)
		}

	}
	return result, err
}

func arrayFromCoords(input []geos.Coord) [][]float64 {
	var result [][]float64
	for inx := 0; inx < len(input); inx++ {
		arr := [...]float64{input[inx].X, input[inx].Y}
		result = append(result, arr[:])
	}
	return result
}

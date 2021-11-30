// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sync"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 700                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)
var waitGroup sync.WaitGroup

type PolygonPerIndex struct {
	Index int
	Polygon string
}

func main() {
	handler := func(
		responseWriter http.ResponseWriter, 
		request *http.Request) {
		get_svg(responseWriter, request)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func getPolygon(i, j, cellIndex int, polygonChannel chan<- PolygonPerIndex) {
	defer waitGroup.Done()
	ax, ay, az := corner(i+1, j)
	bx, by, bz := corner(i, j)
	cx, cy, cz := corner(i, j+1)
	dx, dy, dz := corner(i+1, j+1)
	z_coords := [...]float64{az, bz, cz, dz}
	polygon_color := "#0000ff"
	for _, z := range z_coords {
		if z > 0 {
			polygon_color = "#ff0000"
			break
		}
	}
	polygonString := fmt.Sprintf(
		"<polygon fill=%q points='%g,%g %g,%g %g,%g %g,%g'/>\n",
		polygon_color, ax, ay, bx, by, cx, cy, dx, dy)
	
	polygonChannel <- PolygonPerIndex{cellIndex, polygonString}
}

func get_svg(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "image/svg+xml")
	firstLine := fmt.Sprintf(
		"<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	_, err := io.WriteString(responseWriter, firstLine)
	
	if err != nil {
		fmt.Println("Can not write to response.", err)
		return
	}
	log.Println("Start")
	polygonChannel := make(chan PolygonPerIndex, cells * cells)

	cellIndex := 0
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			waitGroup.Add(1)
			go getPolygon(i, j, cellIndex, polygonChannel)
			cellIndex++
		}
	}
	go func() {
		waitGroup.Wait()
		close(polygonChannel)
	}()
	var allPolygons [cells * cells]string
	for {
		polygon, channelIsOpen := <-polygonChannel
		if !channelIsOpen {
			break
		}
		allPolygons[polygon.Index] = polygon.Polygon
	}
	for _, polygon := range allPolygons {
		_, err = io.WriteString(responseWriter, polygon)
		
		if err != nil {
			fmt.Printf(
				"Can not write the polygon %s. \nReason: %v", 
				polygon, 
				err,
			)
		}
	}
	_, err = io.WriteString(responseWriter, "</svg>")
	if err != nil {
		fmt.Println("Can not write the final tag '</svg>': ", err)
	}
	log.Println("Finish")
}

func corner(i, j int) (projectedX, projectedY, surfaceHeight float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	surfaceHeight = f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	projectedX = width/2 + (x-y)*cos30*xyscale
	projectedY = height/2 + (x+y)*sin30*xyscale - surfaceHeight*zscale
	return
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-

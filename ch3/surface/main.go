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
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	handler := func(
		responseWriter http.ResponseWriter, 
		request *http.Request) {
		get_svg(responseWriter, request)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
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

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
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
			_, err = io.WriteString(responseWriter, polygonString)
			
			if err != nil {
				fmt.Printf(
					"Can not write the polygon %s. \nReason: %v", 
					polygonString, 
					err)
			}
		}
	}
	_, err = io.WriteString(responseWriter, "</svg>")
	if err != nil {
		fmt.Println("Can not write to the file.", err)
	}
}

func corner(i, j int) (float64, float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 61.
//!+

// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", getImage)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getImage(responseWriter http.ResponseWriter, request *http.Request) {
	var height, width int
	var xmax, xmin, ymax, ymin float64
	
	url := request.URL
	urlArgs := url.Query()
	parsedWidth, err := strconv.ParseInt(urlArgs.Get("width"), 0, 0)
	if err != nil {
		width = 1024
	} else {
		width = int(parsedWidth)
	}

	parsedHeight, err := strconv.ParseInt(urlArgs.Get("height"), 0, 0)
	if err != nil {
		height = 1024
	} else {
		height = int(parsedHeight)
	}

	if width <=0 || height <= 0 {
		responseWriter.WriteHeader(400)
		_, _ = io.WriteString(
			responseWriter, 
			"Bad request." + 
			"The width and height should be positive")
		return
	}

	xmin, err = strconv.ParseFloat(urlArgs.Get("xmin"), 64)
	if err != nil {
		xmin = -2
	}

	xmax, err = strconv.ParseFloat(urlArgs.Get("xmax"), 64)
	if err != nil {
		xmax = +2
	}

	ymin, err = strconv.ParseFloat(urlArgs.Get("ymin"), 64)
	if err != nil {
		ymin = -2
	}

	ymax, err = strconv.ParseFloat(urlArgs.Get("ymax"), 64)
	if err != nil {
		ymax = +2
	}

	if xmin >= xmax || ymin >= ymax {
		responseWriter.WriteHeader(400)
		_, _ = io.WriteString(
			responseWriter, 
			"Bad request." + 
			"xmin and ymin should be less than xmax and ymax respectively")
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			x1 := (float64(px) - 0.25)/float64(width)*(xmax-xmin) + xmin
			x2 := (float64(px) + 0.25)/float64(width)*(xmax-xmin) + xmin
			y1 := (float64(py) - 0.25)/float64(height)*(ymax-ymin) + ymin
			y2 := (float64(py) + 0.25)/float64(height)*(ymax-ymin) + ymin
			z1 := complex(x1, y1)
			z2 := complex(x2, y1)
			z3 := complex(x1, y2)
			z4 := complex(x2, y2)
			r1, g1, b1, a1 := newton(z1).RGBA()
			r2, g2, b2, a2  := newton(z2).RGBA()
			r3, g3, b3, a3  := newton(z3).RGBA()
			r4, g4, b4, a4  := newton(z4).RGBA()
			smoothedColor := color.RGBA{
				uint8((r1 + r2 + r3 + r4) / 4),
				uint8((g1 + g2 + g3 + g4) / 4),
				uint8((b1 + b2 + b3 + b4) / 4),
				uint8((a1 + a2 + a3 + a4) / 4),
			}
			img.Set(px, py, smoothedColor)
		}
	}
	png.Encode(responseWriter, img) // NOTE: ignoring errors
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.RGBA{100, 200, 255 - contrast*n, 255}
		}
	}
	return color.Black
}

//!-

// Some other interesting functions:

func acos(z complex128) color.Color {
	v := cmplx.Acos(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}

func sqrt(z complex128) color.Color {
	v := cmplx.Sqrt(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{128, blue, red}
}

// f(x) = x^4 - 1
//
// z' = z - f(z)/f'(z)
//    = z - (z^4 - 1) / (4 * z^3)
//    = z - (z - 1/z^3) / 4
func newton(z complex128) color.Color {
	const iterations = 37
	const contrast = 7
	for i := uint8(0); i < iterations; i++ {
		z -= (z - 1/(z*z*z)) / 4
		if cmplx.Abs(z*z*z*z-1) < 1e-6 {
			return color.Gray{255 - contrast*i}
		}
	}
	return color.Black
}

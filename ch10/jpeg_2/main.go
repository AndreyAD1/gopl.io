// The jpeg command reads a PNG image from the standard input
// and writes it as a JPEG image to the standard output.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png" // register PNG decoder
	"io"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprint(os.Stderr, "expect one output format to be specified\n")
		os.Exit(1)
	}
	outputFormat := flag.Args()[0]
	if outputFormat != "jpeg" && outputFormat != "png" {
		fmt.Fprint(os.Stderr, "an output format should be 'jpeg' or 'png'\n")
		os.Exit(1)
	}
	if err := toJPEG(os.Stdin, os.Stdout, outputFormat); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer, outputFormat string) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	var encoderError error
	switch {
	case outputFormat == "jpeg" :
		encoderError = jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case outputFormat == "png" :
		encoderError = png.Encode(out, img)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
	return encoderError
}

//!-main

/*
//!+with
$ go build gopl.io/ch3/mandelbrot
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
Input format = png
//!-with

//!+without
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
jpeg: image: unknown format
//!-without
*/

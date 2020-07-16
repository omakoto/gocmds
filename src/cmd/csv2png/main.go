package main

import (
	"encoding/csv"
	"github.com/omakoto/go-common/src/common"
	"github.com/pborman/getopt/v2"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	debug   = getopt.BoolLong("debug", 'd', "Enable debug output")
	verbose = getopt.BoolLong("verbose", 'v', "Enable verbose output")

	grayscale = getopt.BoolLong("grayscale", 'g', "Treat input values as float64 and create grayscale image")

	outFile = getopt.StringLong("out", 'o', "out.png", "Output file name")
	inFile  = getopt.StringLong("in", 'i', "-", "Input file name")
)

func generateGrayscaleImage(values [][]string) image.Image {
	h := len(values[0])
	w := len(values)
	image := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})

	var fmatrix [][]float64
	fmatrix = make([][]float64, h)
	var max, v float64
	var err error

	for y := 0; y < h; y++ {
		fmatrix[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			s := strings.TrimSpace(values[x][y])
			if s == "" {
				v = 0
			} else {
				v, err = strconv.ParseFloat(s, 64)
			}
			common.Checkf(err, "unable to parse %#v at row %d, column %d", values[x][y], y, x)
			fmatrix[y][x] = v
			if v > max {
				max = v
			}
		}
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			color := color.Gray{uint8(fmatrix[y][x] / max * 256)}
			image.SetGray(x, y, color)
		}
	}

	return image
}

func realMain() int {
	getopt.Parse()
	if *debug {
		common.DebugEnabled = true
	}
	if *verbose {
		common.VerboseEnabled = true
	}

	var err error
	var in io.ReadCloser
	if *inFile == "-" {
		common.Verbosef("Reading from stdin...")
		in = os.Stdin
	} else {
		common.Verbosef("Reading from %s...", *inFile)
		in, err = os.Open(*inFile)
		common.Checkf(err, "unable to open file")
	}
	defer in.Close()

	// Read from the CSV file.
	rd := csv.NewReader(in)
	values, err := rd.ReadAll()
	common.Checkf(err, "unable to read csv file")

	// Convert to an image.
	var image image.Image
	if *grayscale {
		image = generateGrayscaleImage(values)
	} else {
		panic("RGB mode not implemented yet.")
	}

	// Write to a file.
	out, err := os.OpenFile(*outFile, os.O_WRONLY|os.O_CREATE, 0777)
	common.Checkf(err, "unable to write to file")
	defer out.Close()

	err = png.Encode(out, image)
	common.Checke(err)

	return 0
}

func main() {
	common.RunAndExit(realMain)
}

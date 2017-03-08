package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/adamryman/radialwave/parse"
	"github.com/adamryman/radialwave/viz"
	"github.com/mjibson/go-dsp/wav"
	flag "github.com/spf13/pflag"
)

var (
	r    = flag.IntP("radius", "r", 3000, "radius")
	file = flag.StringP("outfile", "o", "output.png", "file output")
	fill = flag.Float64P("fill", "f", 1, "fill between sections")
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Need a second argument")
	}
	in, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w, err := wav.New(in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(w.SampleRate)
	fmt.Println(w.BitsPerSample)
	freq, err := parse.Wav(w)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f, err := os.Create(*file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	circle, err := viz.Circle(freq, *r, *fill)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = png.Encode(f, circle)
	f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

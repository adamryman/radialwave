package main

import (
	"fmt"
	"image/png"
	"os"
	"strconv"

	"github.com/adamryman/radialwave/parse"
	"github.com/adamryman/radialwave/viz"
	"github.com/mjibson/go-dsp/wav"
	flag "github.com/spf13/pflag"

	"github.com/pkg/errors"
)

var (
	radius = flag.IntP("radius", "r", 3000, "radius")
	fill   = flag.Float64P("fill", "f", 1, "fill between sections")

	// TODO: Implement properly
	simple = flag.Int("simple", 0, "produce simple circle with passed number of chords")

	animate = flag.BoolP("animate", "a", false, "output pngs to be animated")

	outFile = flag.StringP("outfile", "o", "output.png", "file output")
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()
	var errMessage error
	defer func() {
		if errMessage != nil {
			fmt.Println(errMessage)
		}
	}()

	var renderFunc func([]int) error
	if *animate {
		renderFunc = renderPNGs
	} else {
		renderFunc = renderPNG
	}

	if *simple != 0 {
		freq := make([]int, *simple)
		for i := 0; i < len(freq); i++ {
			freq[i] = i
		}

		err := renderFunc(freq)
		if err != nil {
			errMessage = errors.Wrap(err, "error creating simple circle")
			return 1
		}
		return 0
	}

	if len(flag.Args()) < 1 {
		errMessage = errors.New("need input file")
		return 1
	}
	input := flag.Arg(0)

	err := handleInputWav(input, renderFunc)
	if err != nil {
		errMessage = err
		return 1
	}

	return 0
}

func renderPNGs(freq []int) error {
	circles, err := viz.Circles(freq, *radius, *fill)
	if err != nil {
		return err
	}

	for i, v := range circles {
		f, err := os.Create(strconv.Itoa(i) + "_" + *outFile)
		if err != nil {
			return err
		}

		//freq := []int{0, 100, 1000, 100000}

		err = png.Encode(f, v)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil

}

func renderPNG(freq []int) error {
	// TODO: Break this up so I can create a simple cirle with the number of chords passed
	circle, err := viz.Circle(freq, *radius, *fill)
	if err != nil {
		return err
	}

	f, err := os.Create(*outFile)
	if err != nil {
		return err
	}

	//freq := []int{0, 100, 1000, 100000}

	err = png.Encode(f, circle)
	f.Close()
	if err != nil {
		return err
	}

	return nil
}

func handleInputWav(input string, renderFunc func([]int) error) error {
	in, err := os.Open(input)
	if err != nil {
		return errors.Wrapf(err, "cannot open file %s", input)
	}
	w, err := wav.New(in)
	if err != nil {
		return errors.Wrapf(err, "cannot open %s as wav file", input)
	}
	fmt.Println(w.SampleRate)
	fmt.Println(w.BitsPerSample)
	freq, err := parse.WavIntoMaxAmplitudeFrequencies(w)
	if err != nil {
		return errors.Wrapf(err, "cannot parse %s as wav file", input)
	}
	err = renderFunc(freq)
	if err != nil {
		return err
	}

	return nil
}

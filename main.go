package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/adamryman/radialwave/parse"
	"github.com/adamryman/radialwave/viz"
	"github.com/mjibson/go-dsp/wav"
	flag "github.com/spf13/pflag"
	"sync"

	"github.com/pkg/errors"
)

var (
	radius = flag.IntP("radius", "r", 3000, "radius")
	fill   = flag.Float64P("fill", "f", 1, "fill between sections")
	// TODO: Input validation
	bpm = flag.Float64P("bpm", "b", 60, "bpm of input, one arc of the circle for every beat in the input.")

	// TODO: Implement properly
	simple = flag.Int("simple", 0, "produce simple circle with passed number of chords")

	animate = flag.StringP("animate", "a", "", "output pngs to be animated, will put frames into directory passed into flag")
	ffmpeg  = flag.Bool("ffmpeg", false, "run ffmpeg for creating animation with music")

	outFile = flag.StringP("outfile", "o", "output.png", "file output")
)

var inputFile string

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()
	var errMessage error
	defer func() {
		if errMessage != nil {
			fmt.Fprintln(os.Stderr, errMessage)
		}
	}()

	var renderFunc func([]int) error
	if *animate == "" {
		renderFunc = renderPNG
	} else {
		renderFunc = renderPNGs
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
	inputFile = flag.Arg(0)

	err := handleInputWav(inputFile, renderFunc)
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

	var wg sync.WaitGroup

	err = os.MkdirAll(*animate, 0777)
	if err != nil {
		return err
	}
	type indexCircle struct {
		Index  int
		Circle image.Image
	}
	circleChan := make(chan indexCircle)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for c := range circleChan {
				f, err := os.Create(filepath.Join(*animate, strconv.Itoa(c.Index)+".png"))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				err = png.Encode(f, c.Circle)
				f.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}()
	}
	for i, c := range circles {
		circleChan <- indexCircle{i, c}
	}
	close(circleChan)

	wg.Wait()

	bps := *bpm / 60.0
	fmt.Printf("ffmpeg -i %s -r %f -i ./%s/%%d.png -c:v libvpx -b:v 1M -crf 4 %s\n", inputFile, bps, *animate, *outFile)
	if *ffmpeg {
		ffmpegCmd := exec.Command("ffmpeg", strings.Split(fmt.Sprintf("-i %s -r %f -i ./%s/%%d.png -c:v libvpx -b:v 1M -crf 4 %s", inputFile, bps, *animate, *outFile), " ")...)
		ffmpegCmd.Stdin = os.Stdin
		ffmpegCmd.Stderr = os.Stderr
		ffmpegCmd.Stdout = os.Stdout
		err = ffmpegCmd.Run()
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
	freq, err := parse.WavIntoMaxAmplitudeFrequencies(w, *bpm)
	if err != nil {
		return errors.Wrapf(err, "cannot parse %s as wav file", input)
	}
	err = renderFunc(freq)
	if err != nil {
		return err
	}

	return nil
}

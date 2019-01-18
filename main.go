package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
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

	chirp = flag.Int("chirp", 0, "produces a circle using a chrirp with passed number of frequencies/chords")

	animate = flag.StringP("animate", "a", "", "output pngs to be animated, will put frames into directory passed into flag")
	ffmpeg  = flag.Bool("ffmpeg", false, "run ffmpeg for creating animation with music")

	outFile = flag.StringP("outfile", "o", "output.png", "file output")
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	flag.Parse()

	var inputFile string
	var err error

	freq := createChirp(*chirp)

	// If we don't have frequencies at this point, chirp was not used
	if len(freq) == 0 {
		if len(flag.Args()) < 1 {
			return errors.New("need input file")
		}
		inputFile = flag.Arg(0)

		log.Printf("Extracting frequency content of input file %s\n", inputFile)
		freq, err = getMaxFrequenciesFromWav(inputFile)
		if err != nil {
			return err
		}
	}
	animatePath := *animate

	if animatePath != "" {
		circles, err := viz.Circles(freq, *radius, *fill)
		if err != nil {
			return err
		}
		err = renderPNGs(circles, animatePath)
		if err != nil {
			return err
		}

		if *ffmpeg {
			bps := *bpm / 60.0
			// TODO: ffmpeg -i harder.wav -r 2.053333 -i ./harder_highdef/%d.png -c:v libx264 -b:v 1M -f mp4 -preset fast -pix_fmt yuv420p -r 2.053333 -acodec aac harder_highdef.mp4
			// TODO: make mp4
			log.Printf("ffmpeg -i %s -r %f -i ./%s/%%d.png -c:v libvpx -b:v 1M -crf 4 %s\n", inputFile, bps, animatePath, *outFile)
			// TODO: ffmpeg -i harder.wav -r 2.053333 -i ./harder_highdef/%d.png -c:v libx264 -b:v 1M -f mp4 -preset fast -pix_fmt yuv420p -r 2.053333 -acodec aac harder_highdef.mp4
			// TODO: make mp4
			ffmpegCmd := exec.Command("ffmpeg", strings.Split(fmt.Sprintf("-i %s -r %f -i ./%s/%%d.png -c:v libvpx -b:v 1M -crf 4 %s", inputFile, bps, animatePath, *outFile), " ")...)
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

	circle, err := viz.Circle(freq, *radius, *fill)
	if err != nil {
		return err
	}
	err = renderPNG(circle)
	if err != nil {
		return err
	}

	return nil
}

func createChirp(chirpCount int) []int {
	freq := make([]int, chirpCount)
	for i := 0; i < len(freq); i++ {
		freq[i] = i
	}
	return freq
}

func getMaxFrequenciesFromWav(inputFile string) ([]int, error) {
	in, err := os.Open(inputFile)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open file %s", inputFile)
	}
	w, err := wav.New(in)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open %s as wav file", inputFile)
	}
	freq, err := parse.WavIntoMaxAmplitudeFrequencies(w, *bpm)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse %s as wav file", inputFile)
	}
	return freq, nil
}

// renderPNGs takes a slice of images, and renders them into the directory path
func renderPNGs(pngs []image.Image, path string) error {

	var wg sync.WaitGroup

	// create directory
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	// create indexable pngs for file names
	type indexedPNG struct {
		Index int
		PNG   image.Image
	}

	// create channel for work to be added to
	pngChan := make(chan indexedPNG)
	// put images into channel to be processed
	go func() {
		for i, p := range pngs {
			log.Printf("Rending frame %d / %d\n", i, len(pngs))
			pngChan <- indexedPNG{i, p}
		}
		close(pngChan)
	}()

	// Create workers
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range pngChan {
				f, err := os.Create(filepath.Join(path, strconv.Itoa(p.Index)+".png"))
				if err != nil {
					log.Println(err)
				}

				err = png.Encode(f, p.PNG)
				f.Close()
				if err != nil {
					log.Println(err)
				}
			}
		}()
	}

	// Wait for workers to finish
	wg.Wait()

	return nil
}

func renderPNG(circle image.Image) error {
	f, err := os.Create(*outFile)
	if err != nil {
		return err
	}

	err = png.Encode(f, circle)
	f.Close()
	if err != nil {
		return err
	}

	return nil
}

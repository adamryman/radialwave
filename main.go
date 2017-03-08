package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/adamryman/visualizer-experiment/parse"
	"github.com/adamryman/visualizer-experiment/viz"
	"github.com/mjibson/go-dsp/wav"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need a second argument")
	}
	in, err := os.Open(os.Args[1])
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
	f, err := os.Create("output.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	circle, err := viz.Circle(freq)
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

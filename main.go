package main

import (
	"os"

	"github.com/adamryman/radialwave/viz"
	"github.com/mjibson/go-dsp/wav"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need a second argument")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w, err := wav.New(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(w.SampleRate)
	fmt.Println(w.BitsPerSample)
	err := Circle(w)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

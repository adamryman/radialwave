// This is very unfinished, but I wanted to share
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/y0ssar1an/q"

	"github.com/adamryman/radialwave/wave"
)

// TODO:
// - read in data chunks
// - possibly assign chunks numbers to recombine later
// - pass groups of chunks to something to do fft on them.
// - Find dominate frequency of group of chunks
// - Pass dominate frequency info to both something to create a slice of the frequency data in
//   original music AND something to create a sorted unique slice of frequencies,
//   this could be done effectily with insertion sort
// - After all pieces have been collected:
// - Create mapping from wavelength to color
// - Render some somehow

func main() {
	b, err := ioutil.ReadFile("./harder.wav")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := bytes.NewReader(b)

	w := wave.Format{ChunkID: 1}
	original := r.Len()
	Q(original)
	err = w.Decode(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Q(w)
	Q(r.Len())
	Q(original - r.Len())
	Q(len(w.Data))

	Q(r.Len())
	err = w.Decode(r)
	if err != nil {
		Q("COULD NOT DECODE AGAIN")
		Q(r.Len())
		Q(err)
		fmt.Println(err)
		os.Exit(1)

	}
	Q(w)
	Q(r.Len())
	b2, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Q(len(b2))

	//err = w.Decode(r)
	//if err != nil {
	//fmt.Println(err)
	//os.Exit(1)
	//}
	//Q(w)
	//Q(original - r.Len())
}

type chunk struct {
	id      int
	rawData []byte
	fftData []float64
}

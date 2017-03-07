package parse

import (
	"fmt"
	"log"
	"math/cmplx"
	"os"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
	. "github.com/y0ssar1an/q"
)

// harder.wav is 226 seconds long
// measuring 452 bits

var logger = log.New(os.Stderr, "Render", log.LstdFlags)

func Wav(w *wav.Wav) ([]int, error) {
	Q(int(w.BitsPerSample))
	Q(w.Samples)
	Q(int(w.SampleRate))
	fmt.Println(w.NumChannels)
	seconds := (float64(w.Samples) / float64(w.SampleRate))
	_ = seconds
	//for err == nil {
	//samples, err = w.ReadSamples(int(w.SampleRate))
	samples, err := w.ReadFloats(int(w.SampleRate) * int(w.NumChannels))
	_ = err
	l, r := SplitChannels(samples)

	var maxFrequencies []int

	for err == nil {
		fftOutL := fft.FFTReal(float32To64(l))
		fftOutR := fft.FFTReal(float32To64(r))

		maxIndexL, maxL := maxFrequency(fftOutL)
		maxIndexR, maxR := maxFrequency(fftOutR)
		if maxL > maxR {
			maxFrequencies = append(maxFrequencies, maxIndexL)
		} else {
			maxFrequencies = append(maxFrequencies, maxIndexR)
		}

		samples, err = w.ReadFloats(int(w.SampleRate) * int(w.NumChannels))
		l, r = SplitChannels(samples)
	}
	Q(err)
	Q(maxFrequencies)

	return maxFrequencies, nil
}

func SplitChannels(stereo []float32) (l, r []float32) {
	for i := 0; i < len(stereo)/2; i = i + 2 {
		l = append(l, stereo[i])
		r = append(r, stereo[i])
	}

	return l, r
}

func maxFrequency(f []complex128) (int, float64) {
	var maxIndex int
	var max float64

	for i, v := range f {
		if i == 0 {
			continue
		}
		if i > len(f)/2 {
			break
		}
		value := cmplx.Abs(v)
		if value > max {
			max = value
			maxIndex = i
		}
	}
	//Q(maxIndex)

	return maxIndex, max
}

func float32To64(f []float32) []float64 {
	o := make([]float64, len(f))
	for i, v := range f {
		o[i] = float64(v)
	}
	return o
}

// MVP
// Take in song
// Do FFT on each Second

// Image process

// Input
// Music | Time | Frequencies in each unit of time

// Output
// Image | Space | Color

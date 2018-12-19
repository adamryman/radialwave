package parse

import (
	"io"
	"log"
	"math/cmplx"
	"os"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
	. "github.com/y0ssar1an/q"

	"github.com/pkg/errors"
)

// harder.wav is 226 seconds long
// measuring 452 bits

var logger = log.New(os.Stderr, "Render", log.LstdFlags)

// TODO
//func WavIntoMaxAmplitudeFrequenciesByTime(w *wav.Wav, t time.Time) ([]int, error) {

// Wav
// TODO: Fix this not making sense if there is mono audio, or anything other than two channels
func WavIntoMaxAmplitudeFrequencies(w *wav.Wav) ([]int, error) {
	Q(int(w.BitsPerSample))
	Q(w.Samples)
	Q(int(w.SampleRate))

	samplesPerSecond := int(w.SampleRate) * int(w.NumChannels)

	var maxFrequencies []int
	var err error
	var samples []float32
	for err == nil {
		// Read one second of sound
		samples, err = w.ReadFloats(samplesPerSecond)

		// TODO: Update this to handle any number of channels, rather than hard coded to two. Woops. My bad
		l, r := SplitChannels(samples)
		fftOutL := fft.FFTReal(float32To64(l))
		fftOutR := fft.FFTReal(float32To64(r))

		maxIndexL, maxL := maxFrequency(fftOutL)
		maxIndexR, maxR := maxFrequency(fftOutR)
		if maxL > maxR {
			maxFrequencies = append(maxFrequencies, maxIndexL)
		} else {
			maxFrequencies = append(maxFrequencies, maxIndexR)
		}
	}
	if err != nil && err != io.EOF {
		// TODO: Why is an unexpected EOF sometimes happening
		if err != io.ErrUnexpectedEOF {
			return nil, errors.Wrap(err, "error processing sample")
		}
	}

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

package viz

import (
	"log"
	"os"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
)

var logger = log.New(os.Stderr, "Render", log.LstdFlags)

func Circle(w *wav.Wav, chordNum int) error {
	seconds := w.Samples / int(w.SampleRate)

	samplesPerChord := chordNum * w.Samples
	frequencies := make([]float64, chordNum)
	uniquesChan := make(chan float64)

	_ = frequencies

	// We stop early and use the rest of the data for the lasts chord
	for i := 0; i < seconds; i++ {
		nextWav, err := w.ReadFloats(int(w.SampleRate))
		if err != nil {

		}
	}
	return nil
}

// Upsample

func processSample(sample []float32) (float64, error) {
	ya = fft.FFTReal(float32To64(sample))

	return 0, nil
}

func float32To64(in []float32) []float64 {
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = float64(v)
	}
	return out
}

func sortUniques(uniques []float64, in chan float64) {
	next := <-in
	for i, v := range uniques {
		if next >= v {

		}

	}
}

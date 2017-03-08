package viz

import (
	"image"
	"image/color"
	//"math"

	"github.com/adamryman/circle"
	"github.com/adamryman/visualizer-experiment/wavelength"

	. "github.com/y0ssar1an/q"
)

// Draw a circle on a 1000x1000 pixel background
func Circle(freq []int, radius int, fill float64) (image.Image, error) {
	min, max := minMax(freq)
	top := max - min

	var cs []color.Color
	for _, f := range freq {
		dist := float32(f) / float32(top)
		Q(dist)
		// 265 is a magic number
		wlength := int(dist * float32(265))
		c := wavelength.WaveToRGB(wlength)
		cs = append(cs, c)
	}
	img := circle.ColorCircle(radius, fill, cs...)

	return img, nil
}

func minMax(data []int) (min, max int) {
	for _, v := range data {
		if v < min {
			min = v
			continue
		} else if v > max {
			max = v
			continue
		}
	}
	return
}

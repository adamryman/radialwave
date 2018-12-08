package viz

import (
	"image"
	"image/color"
	"math"

	"github.com/adamryman/radialwave/circle"
	"github.com/adamryman/radialwave/wavelength"

	. "github.com/y0ssar1an/q"
)

// Draw a circle on a 1000x1000 pixel background
//func Circle1(freq []int, radius int, fill float64) (image.Image, error) {
//min, max := minMax(freq)
//top := max - min

//var cs []color.Color
//for _, f := range freq {
//// ratio of current frequency
//dist := float32(f) / float32(top)
//// 265 is a magic number
//wlength := int(dist * float32(265))
//c := wavelength.WaveToRGB(wlength)
//cs = append(cs, c)
//}
//Q(cs)
//img := circle.ColorCircle(radius, fill, cs...)

//return img, nil
//}

// Draw a circle on a 1000x1000 pixel background
func Circle(freq []int, radius int, fill float64) (image.Image, error) {
	min, max := minMax(freq)
	top := max - min
	Q(top)
	Q(max)
	Q(min)
	// If min is zero, we should shift everything up by 1. As 0 will map to zero and 1 will map to zero
	// TODO: Figure out how to space numbers equaily

	toplog := math.Log(float64(top))
	if min == 0 {
		toplog = math.Log(float64(top + 1))
	}

	var cs []color.Color
	for _, f := range freq {
		if min == 0 {
			f = f + 1
		}
		Q(f)
		dist := float64(0)
		if f-min != 0 {
			dist = math.Log(float64(f-min)) / toplog
		}
		Q(dist)
		// 265 is a magic number
		// the length - 1 of wavelength.ToRGB
		wlength := int(dist * float64(265))
		Q(wlength)
		c := wavelength.WaveToRGB(wlength)
		cs = append(cs, c)
	}
	//Q(cs)
	img := circle.ColorCircle(radius, fill, cs...)

	return img, nil
}

func colorPicker() {
	// TODO: Given best case range of values from 0 to 255, spread evenly
	// TODO: Given a range larger than 0 to 255, make sure each number only resolves once.
	// i.e. range of 0 to 511 then both 0 and 1 would map to 0, we want to spread them out as much as possible, so if there was no 2 or 3, would would map 0 to 0 and 1 to 1.
	// Probably create a mapping and then move if needed?

	//[][]int [[0,1]. [], [4], [6]] => [[0], [1], [4], [6]]

	//[][]int [[0,1]. [2], [5]. []] => [[0], [1], [2], [5]]

	//[][]int [[0,1]. [2], [4,5]. [6]] => [[0,1], [2], [4,5], [6]]

	// Na this is all silly
	// TODO: Sort frequencies
	// TODO: Dedupe? Yeah
	// TODO: map over one at a time and shift by one if needed to keep the colors different

}

func SpectrumCircle(radius int, fill float64) (image.Image, error) {
	var cs []color.Color
	for _, c := range wavelength.ToRGB {
		cs = append(cs, c)
	}

	return circle.ColorCircle(radius, fill, cs...), nil
}

func minMax(data []int) (min, max int) {
	min = math.MaxInt64
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

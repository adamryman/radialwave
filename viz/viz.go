package viz

import (
	"image"
	"image/color"
	"math"
	"sort"

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
	sortedFreq := make([]int, len(freq), len(freq))
	copy(sortedFreq, freq)
	mapper := spreaderSmallToLarge(sortedFreq, len(wavelength.ToRGB)-1, min, max)
	var cs []color.Color
	for _, f := range freq {
		bucket := mapper[f]
		//Q(bucket)
		c := wavelength.WaveToRGB(bucket)
		cs = append(cs, c)
	}
	//Q(cs)
	img := circle.ColorCircle(radius, fill, cs...)

	return img, nil
}

func Circles(freq []int, radius int, fill float64) ([]image.Image, error) {
	min, max := minMax(freq)
	sortedFreq := make([]int, len(freq), len(freq))
	copy(sortedFreq, freq)
	mapper := spreaderSmallToLarge(sortedFreq, len(wavelength.ToRGB)-1, min, max)
	var cs []color.Color
	for _, f := range freq {
		bucket := mapper[f]
		//Q(bucket)
		c := wavelength.WaveToRGB(bucket)
		cs = append(cs, c)
	}
	//Q(cs)
	img := circle.ColorCircles(radius, fill, cs...)

	return img, nil
}

func bucketorZero(in []int, buckets int, max int) map[int]int {
	out := make(map[int]int)
	floatMax := float64(max)
	if floatMax == 0 {
		return out
	}
	for _, v := range in {
		out[v] = int(math.Floor((float64(v) / floatMax) * float64(buckets)))
	}
	return out
}

func bucketor(in []int, buckets int, min, max int) map[int]int {
	out := make(map[int]int)
	bucketRange := float64(max - min)
	if bucketRange == 0 {
		return out
	}
	for _, v := range in {
		out[v] = int(math.Floor((float64(v-min) / bucketRange) * float64(buckets)))
	}
	return out
}

// spread values out for more unique colors
func spreaderLargeToSmall(in []int, buckets int, min, max int) map[int]int {
	out := make(map[int]int)
	dedupOutput := make(map[int]bool)
	bucketRange := float64(max - min)
	Q(bucketRange)
	if bucketRange == 0 {
		return out
	}
	sort.Ints(in)
	Q(len(in))
	// range from top to bottom
	for i := len(in) - 1; i > 0; i-- {
		v := in[i]
		Q(v)
		// dedupe existing values
		if _, ok := out[v]; ok {
			Q("continue")
			continue
		}
		// map
		out[v] = int(math.Floor((float64(v-min) / bucketRange) * float64(buckets)))
		// If we already seen the output, shift it down by 1
		if dedupOutput[out[v]] && out[v] > 0 {
			out[v] = out[v] - 1
		}
		dedupOutput[out[v]] = true
	}
	return out
}

// spread values out for more unique colors
func spreaderSmallToLarge(in []int, buckets int, min, max int) map[int]int {
	out := make(map[int]int)
	dedupOutput := make(map[int]bool)
	bucketRange := float64(max - min)
	Q(bucketRange)
	if bucketRange == 0 {
		return out
	}
	sort.Ints(in)
	Q(len(in))
	// range from bottom to top
	for _, v := range in {
		// dedupe existing values
		if _, ok := out[v]; ok {
			Q("continue")
			continue
		}
		// map
		out[v] = int(math.Floor((float64(v-min) / bucketRange) * float64(buckets)))

		// TODO: How to invert colors
		//out[v] = int(math.Floor((math.Abs((1 - (float64(v-min) / bucketRange))) * float64(buckets))))

		// If we already seen the output, shift it down by 1
		if dedupOutput[out[v]] && out[v] < buckets {
			out[v] = out[v] + 1
		}
		dedupOutput[out[v]] = true
	}
	return out
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

package circle

import (
	//. "github.com/y0ssar1an/q"
	"image"
	"image/color"
	"math"
)

func Simple(x, y, radius int) *Circle {
	return &Circle{
		Point:  image.Pt(x, y),
		Radius: radius,
	}
}

// Circle implements image.Image
// godoc image Image
type Circle struct {
	Point  image.Point
	Radius int
}

func (c *Circle) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(
		c.Point.X-c.Radius, // x0
		c.Point.Y-c.Radius, // y0
		c.Point.X+c.Radius, // x1
		c.Point.Y+c.Radius, // y1
	)
}

func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.Point.X)+0.5, float64(y-c.Point.Y)+0.5, float64(c.Radius)
	if xx*xx+yy*yy < rr*rr {
		return color.Black
	}
	return color.White
}

// Sector coming soon
type Sector struct {
	Circle
	Θ1    float64
	Θ2    float64
	Color color.Color
}

func (s *Sector) At(x, y int) color.Color {
	// Center around origin
	xx, yy, rr := float64(x-s.Point.X), -float64(y-s.Point.Y), float64(s.Radius)
	if xx*xx+yy*yy < rr*rr {
		theta := math.Atan2(yy, xx)
		if theta < 0 {
			theta = theta + 2*math.Pi
		}
		if theta >= s.Θ1 && theta < s.Θ2 {
			return s.Color
		}
	}
	return color.Alpha{255}
}

type SectorCircle struct {
	Circle
	Sectors []Sector
}

func (s *SectorCircle) At(x, y int) color.Color {
	// Center around origin
	xx, yy, rr := float64(x-s.Point.X), -float64(y-s.Point.Y), float64(s.Radius)
	if xx*xx+yy*yy < rr*rr {
		theta := math.Atan2(yy, xx)
		if theta < 0 {
			theta = theta + 2*math.Pi
		}
		distance := theta / (2 * math.Pi)
		sectorIndex := int(distance * float64(len(s.Sectors)))
		return s.Sectors[sectorIndex].Color
	}
	return color.Alpha{255}
}

type colorCircle struct {
	Circle
	Colors []color.Color
	Fill   float64
}

func (c *colorCircle) At(x, y int) color.Color {
	// Center around origin
	xx, yy, rr := float64(x-c.Point.X), -float64(y-c.Point.Y), float64(c.Radius)
	if xx*xx+yy*yy < rr*rr {
		theta := math.Atan2(yy, xx)
		if theta < 0 {
			theta = theta + 2*math.Pi
		}
		distance := theta / (2 * math.Pi)
		colorIndex := distance * float64(len(c.Colors))
		diff := colorIndex - math.Floor(colorIndex)

		// 1 = filled, .5 = half filled. 0 = no fill
		if diff > c.Fill {
			return color.Alpha{255}
		}
		return c.Colors[int(colorIndex)]
	}
	return color.Alpha{255}
}

func ColorCircle(radius int, fill float64, colors ...color.Color) image.Image {
	if fill == 0 || fill > 1 {
		fill = 1
	}
	cc := colorCircle{
		Circle: Circle{
			Radius: radius,
			Point:  image.Pt(0, 0),
		},
		Fill: fill,
	}
	for _, c := range colors {
		cc.Colors = append(cc.Colors, c)
	}
	return &cc
}

type Red struct{}

func (_ Red) RGBA() (r, g, b, a uint32) {
	return 0xff, 0, 0, 0xff
}

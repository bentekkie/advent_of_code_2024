package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"github.com/billglover/gradient"
)

var gifPath = flag.String("gif", "", "")

func main() {
	flag.Parse()
	input := inputs.String()
	part1(input)
	part2(input)
}

type P struct {
	x, y int
}

func (p P) up() P {
	return P{p.x, p.y - 1}
}

func (p P) down() P {
	return P{p.x, p.y + 1}
}

func (p P) left() P {
	return P{p.x - 1, p.y}
}

func (p P) right() P {
	return P{p.x + 1, p.y}
}
func (p P) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}
func (p P) inbounds(minx, miny, maxx, maxy int) bool {
	return p.x >= minx && p.x <= maxx && p.y >= miny && p.y <= maxy
}

type D int

const (
	up    D = 0
	right D = 1
	down  D = 2
	left  D = 3
)

func part1(input string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var guard P
	guarDir := up
	obsticals := map[P]struct{}{}
	for y, line := range lines {
		line := strings.TrimSpace(line)
		for x, c := range line {
			if c == '#' {
				obsticals[P{x, y}] = struct{}{}
			}
			if c == '^' {
				guard = P{x, y}
			}
		}
	}
	visited := map[P]struct{}{}
	imgs := []*image.Paletted{drawBaseGrid(guard, obsticals, len(lines[0])-1, len(lines)-1)}
	i := 0
	for guard.inbounds(0, 0, len(lines[0])-1, len(lines)-1) {
		visited[guard] = struct{}{}
		var nextP P
		switch guarDir {
		case up:
			nextP = guard.up()
		case down:
			nextP = guard.down()
		case left:
			nextP = guard.left()
		case right:
			nextP = guard.right()
		}
		if _, ok := obsticals[nextP]; ok {
			guarDir = (guarDir + 1) % 4
		} else {
			guard = nextP
		}
		if *gifPath != "" {
			imgs = append(imgs, drawGuard(imgs[len(imgs)-1], guard, g[i%len(g)]))
		}
		i++
	}
	if *gifPath != "" {
		imgs = skipFrames(imgs, 5)
		imf, err := os.Create(*gifPath + ".png")
		if err != nil {
			panic(err)
		}
		defer imf.Close()
		png.Encode(imf, imgs[len(imgs)-1])
		f, err := os.Create(*gifPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		delays := []int{}
		for range len(imgs) {
			delays = append(delays, 2)
		}
		gif.EncodeAll(f, &gif.GIF{
			Image: imgs,
			Delay: delays,
		})
	}
	fmt.Printf("Part 1: %d\n", len(visited))
}

func skipFrames(imgs []*image.Paletted, skip int) []*image.Paletted {
	out := []*image.Paletted{}
	for i, img := range imgs {
		if i%skip == 0 || i == len(imgs)-1 {
			out = append(out, img)
		}
	}
	return out
}

var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}
var blue = color.RGBA{0, 0, 255, 255}
var yellow = color.RGBA{255, 255, 0, 255}
var p = color.Palette{color.White, red, green, blue}
var g []color.Color

func init() {
	for i := range 20 {
		g = append(g, gradient.LinearPoint(blue, yellow, float64(i)/20))
	}
	p = append(p, g...)
	revG := make([]color.Color, len(g))
	for i, c := range g {
		revG[len(g)-1-i] = c
	}
	g = append(g, revG...)
}

func rectForP(p P, size int) image.Rectangle {
	return image.Rect(p.x*size, p.y*size, (p.x+1)*size, (p.y+1)*size)
}

func drawGuard(base *image.Paletted, guard P, c color.Color) *image.Paletted {
	size := 5
	img := &image.Paletted{
		Pix:     make([]uint8, len(base.Pix)),
		Stride:  base.Stride,
		Rect:    base.Rect,
		Palette: base.Palette,
	}
	copy(img.Pix, base.Pix)
	draw.Draw(img,
		rectForP(guard, size),
		&image.Uniform{c},
		image.Point{},
		draw.Over,
	)
	return img
}

func drawBaseGrid(guard P, obsticals map[P]struct{}, maxx, maxy int) *image.Paletted {
	size := 5
	img := image.NewPaletted(image.Rect(0, 0, maxx*size, maxy*size), p)
	draw.Draw(img,
		rectForP(guard, size),
		&image.Uniform{green},
		image.Point{},
		draw.Over,
	)
	for p := range obsticals {
		draw.Draw(img,
			rectForP(p, size),
			&image.Uniform{red},
			image.Point{},
			draw.Over,
		)
	}
	return img
}

type V struct {
	p P
	d D
}

func step(guard V, obsticals map[P]struct{}) V {
	var nextP P
	switch guard.d {
	case up:
		nextP = guard.p.up()
	case down:
		nextP = guard.p.down()
	case left:
		nextP = guard.p.left()
	case right:
		nextP = guard.p.right()
	}
	d := guard.d
	p := guard.p
	if _, ok := obsticals[nextP]; ok {
		d = (d + 1) % 4
	} else {
		p = nextP
	}

	return V{p, d}
}

func path(start V, obsticals map[P]struct{}, maxx, maxy int) map[P]struct{} {
	guard := start.p
	guarDir := start.d
	visited := map[P]struct{}{}
	for guard.inbounds(0, 0, maxx, maxy) {
		visited[guard] = struct{}{}
		var nextP P
		switch guarDir {
		case up:
			nextP = guard.up()
		case down:
			nextP = guard.down()
		case left:
			nextP = guard.left()
		case right:
			nextP = guard.right()
		}
		if _, ok := obsticals[nextP]; ok {
			guarDir = (guarDir + 1) % 4
		} else {
			guard = nextP
		}
	}
	return visited
}

func hasCycle(guard V, obsticals map[P]struct{}, maxx, maxy int) bool {
	t, h := guard, guard
	for h.p.inbounds(0, 0, maxx, maxy) {
		t = step(t, obsticals)
		h = step(step(h, obsticals), obsticals)
		if t == h {
			return true
		}
	}
	return false
}

func part2(input string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var guard P
	guarDir := up
	obsticals := map[P]struct{}{}
	possibleObsticals := []P{}
	for y, line := range lines {
		line := strings.TrimSpace(line)
		for x, c := range line {
			if c == '#' {
				obsticals[P{x, y}] = struct{}{}
			}
			if c == '^' {
				guard = P{x, y}
			}
			if c == '.' {
				possibleObsticals = append(possibleObsticals, P{x, y})
			}
		}
	}
	original := path(V{guard, guarDir}, obsticals, len(lines[0])-1, len(lines)-1)
	var s int
	benlog.Timed(func() {
		for _, o := range possibleObsticals {
			if _, ok := original[o]; !ok {
				continue
			}
			obsticals[o] = struct{}{}
			if hasCycle(V{guard, guarDir}, obsticals, len(lines[0])-1, len(lines)-1) {
				s++
			}
			delete(obsticals, o)
		}
	})
	fmt.Printf("Part 2: %d\n", s)
}

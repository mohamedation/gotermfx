package animations

import (
	"context"
	"gotermfx/termfx"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type starfield struct{}

type star struct {
	x, y, z float64
}

func (s *starfield) Run(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cols, rows := termfx.GetSize()
	numStars := cols * rows / 4

	stars := make([]star, numStars)
	for i := range stars {
		stars[i] = randomStar(rng)
	}

	var buf strings.Builder
	buf.Grow(cols * rows * 8)

	type cell struct {
		r     rune
		color string
	}
	grid := make([]cell, cols*rows)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		newCols, newRows := termfx.GetSize()
		if newCols != cols || newRows != rows {
			cols, rows = newCols, newRows
			numStars = cols * rows / 4
			stars = make([]star, numStars)
			for i := range stars {
				stars[i] = randomStar(rng)
			}
			grid = make([]cell, cols*rows)
		}

		for i := range grid {
			grid[i].r = 0
		}

		cx := float64(cols) / 2.0
		cy := float64(rows) / 2.0

		for i := range stars {
			st := &stars[i]
			st.z -= 0.015
			if st.z <= 0 {
				stars[i] = randomStar(rng)
				continue
			}

			sx := int(st.x/st.z*cx + cx)
			sy := int(st.y/st.z*cy*0.5 + cy)
			if sx < 0 || sx >= cols || sy < 0 || sy >= rows {
				stars[i] = randomStar(rng)
				continue
			}

			brightness := 1.0 - st.z
			var ch rune
			var color string
			switch {
			case brightness > 0.75:
				ch = '*'
				color = "\033[97m" // bright white
			case brightness > 0.5:
				ch = '+'
				color = "\033[37m" // white
			case brightness > 0.25:
				ch = '.'
				color = "\033[90m" // dark gray
			default:
				ch = '·'
				color = "\033[90m"
			}
			idx := sy*cols + sx
			grid[idx].r = ch
			grid[idx].color = color
		}

		buf.Reset()
		buf.WriteString("\033[H")
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				c := grid[i*cols+j]
				if c.r != 0 {
					buf.WriteString(c.color)
					buf.WriteRune(c.r)
					buf.WriteString("\033[0m")
				} else {
					buf.WriteByte(' ')
				}
			}
			if i < rows-1 {
				buf.WriteString("\r\n")
			}
		}
		os.Stdout.WriteString(buf.String())
		time.Sleep(33 * time.Millisecond)
	}
}

func randomStar(rng *rand.Rand) star {
	angle := rng.Float64() * 2 * math.Pi
	radius := rng.Float64()*0.8 + 0.2
	return star{
		x: math.Cos(angle) * radius,
		y: math.Sin(angle) * radius,
		z: 0.5 + rng.Float64()*0.5,
	}
}

func init() {
	termfx.Register("Starfield", &starfield{})
}

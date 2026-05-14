package animations

import (
	"context"
	"github.com/mohamedation/gotermfx/termfx"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type hyperspace struct{}

type swStar struct {
	angle   float64
	r       float64
	speed   float64
	tailLen float64
}

func hypLineChar(angle float64) rune {
	c := math.Cos(angle)
	s := math.Sin(angle)
	ac := math.Abs(c)
	as := math.Abs(s) * 2.0
	if ac > as*1.6 {
		return '-'
	}
	if as > ac*1.6 {
		return '|'
	}
	if c*s >= 0 {
		return '\\'
	}
	return '/'
}

func hypStreakColor(speed float64) string {
	switch {
	case speed < 0.4:
		return "\033[90m" // dark gray
	case speed < 0.9:
		return "\033[37m" // gray
	case speed < 1.8:
		return "\033[97m" // white
	default:
		return "\033[97;1m" // bold white
	}
}

type hypCell struct {
	r     rune
	color string
}

func hypDrawStreak(grid []hypCell, cols, rows int, cx, cy, aspect float64, s *swStar, ch rune, color string) {
	end := s.tailLen
	if end > s.r {
		end = s.r
	}
	const step = 0.55
	for t := 0.0; t <= end; t += step {
		r := s.r - t
		sx := int(cx + math.Cos(s.angle)*r)
		sy := int(cy + math.Sin(s.angle)*r*aspect)
		if sx >= 0 && sx < cols && sy >= 0 && sy < rows {
			idx := sy*cols + sx
			grid[idx].r = ch
			grid[idx].color = color
		}
	}
}

func hypRenderGrid(grid []hypCell, cols, rows int, buf *strings.Builder) {
	buf.Reset()
	buf.WriteString("\033[H")
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			c := grid[i*cols+j]
			if c.r != ' ' {
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
}

func newBuildupStar(rng *rand.Rand) swStar {
	return swStar{
		angle:   rng.Float64() * 2 * math.Pi,
		r:       rng.Float64() * 1.5,
		speed:   0.15 + rng.Float64()*0.3,
		tailLen: 0.2,
	}
}

func newHyperspaceStar(rng *rand.Rand) swStar {
	return swStar{
		angle:   rng.Float64() * 2 * math.Pi,
		r:       rng.Float64() * 4.0,
		speed:   10.0 + rng.Float64()*8.0,
		tailLen: 3.0,
	}
}

func (h *hyperspace) Run(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cols, rows := termfx.GetSize()

	numStars := 90
	phase := 0
	phaseTick := 0

	stars := make([]swStar, numStars)
	for i := range stars {
		stars[i] = newBuildupStar(rng)
		stars[i].r = rng.Float64() * 8.0
	}

	grid := make([]hypCell, cols*rows)
	var buf strings.Builder
	buf.Grow(cols * rows * 8)

	const buildupTicks = 200

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		newCols, newRows := termfx.GetSize()
		if newCols != cols || newRows != rows {
			cols, rows = newCols, newRows
			grid = make([]hypCell, cols*rows)
			buf.Reset()
			buf.Grow(cols * rows * 8)
			for i := range stars {
				stars[i] = newBuildupStar(rng)
			}
			phase = 0
			phaseTick = 0
		}

		cx := float64(cols) / 2.0
		cy := float64(rows) / 2.0
		aspect := cy / cx * 2.0

		phaseTick++

		switch phase {

		case 0:
			for i := range grid {
				grid[i].r = ' '
				grid[i].color = ""
			}

			for i := range stars {
				s := &stars[i]

				s.speed *= 1.013
				if s.speed > 2.5 {
					s.speed = 2.5
				}
				s.r += s.speed

				targetTail := s.speed * 5.0
				if s.tailLen < targetTail {
					s.tailLen += 0.25
				}

				hx := cx + math.Cos(s.angle)*s.r
				hy := cy + math.Sin(s.angle)*s.r*aspect
				if hx < -4 || hx >= float64(cols)+4 || hy < -4 || hy >= float64(rows)+4 {
					stars[i] = newBuildupStar(rng)
					continue
				}

				hypDrawStreak(grid, cols, rows, cx, cy, aspect, s,
					hypLineChar(s.angle), hypStreakColor(s.speed))
			}

			hypRenderGrid(grid, cols, rows, &buf)

			if phaseTick >= buildupTicks {
				phase = 1
				phaseTick = 0
			}

			time.Sleep(35 * time.Millisecond)

		case 1:
			for i := range grid {
				grid[i].r = ' '
				grid[i].color = ""
			}

			for i := range stars {
				s := &stars[i]

				s.speed *= 1.10
				s.r += s.speed

				s.tailLen = s.r * 0.85

				hypDrawStreak(grid, cols, rows, cx, cy, aspect, s,
					hypLineChar(s.angle), "\033[97;1m")
			}

			hypRenderGrid(grid, cols, rows, &buf)

			if phaseTick >= 22 {
				phase = 2
				phaseTick = 0
			}

			time.Sleep(25 * time.Millisecond)

		case 2:
			buf.Reset()
			buf.WriteString("\033[H\033[97m")
			line := strings.Repeat("█", cols)
			for i := 0; i < rows; i++ {
				buf.WriteString(line)
				if i < rows-1 {
					buf.WriteString("\r\n")
				}
			}
			buf.WriteString("\033[0m")
			os.Stdout.WriteString(buf.String())

			if phaseTick >= 4 {
				phase = 3
				phaseTick = 0
				// Seed hyperspace stars.
				for i := range stars {
					stars[i] = newHyperspaceStar(rng)
					stars[i].r = rng.Float64() * 15.0
				}
			}

			time.Sleep(50 * time.Millisecond)

		case 3:
			for i := range grid {
				grid[i].r = ' '
				grid[i].color = ""
			}

			for i := range stars {
				s := &stars[i]
				s.r += s.speed
				s.tailLen = s.r * 0.88

				hx := cx + math.Cos(s.angle)*s.r
				hy := cy + math.Sin(s.angle)*s.r*aspect
				if hx < -2 || hx >= float64(cols)+2 || hy < -2 || hy >= float64(rows)+2 {
					stars[i] = newHyperspaceStar(rng)
					stars[i].r = rng.Float64() * 3.0
					continue
				}

				var color string
				switch rng.Intn(4) {
				case 0:
					color = "\033[97;1m" // bright white
				case 1:
					color = "\033[96;1m" // bright cyan
				case 2:
					color = "\033[94m" // blue
				default:
					color = "\033[97m" // white
				}
				hypDrawStreak(grid, cols, rows, cx, cy, aspect, s,
					hypLineChar(s.angle), color)
			}

			hypRenderGrid(grid, cols, rows, &buf)

			if phaseTick >= 136 {
				if termfx.IsOnce(ctx) {
					return
				}
				phase = 0
				phaseTick = 0
				for i := range stars {
					stars[i] = newBuildupStar(rng)
					stars[i].r = rng.Float64() * 8.0
				}
			}

			time.Sleep(22 * time.Millisecond)
		}
	}
}

func init() {
	termfx.Register("Hyperspace", &hyperspace{})
}

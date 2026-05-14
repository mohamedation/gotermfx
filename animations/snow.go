package animations

import (
	"context"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mohamedation/gotermfx/termfx"
)

type snow struct{}

type snowflake struct {
	x, y  float64
	speed float64
	drift float64
	char  rune
}

var snowChars = []rune{'*', '❄', '·', '•', '+'}

func (s *snow) Run(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cols, rows := termfx.GetSize()

	flakes := make([]snowflake, cols/3)
	for i := range flakes {
		flakes[i] = newFlake(rng, cols, rows, true)
	}

	type cell struct{ r rune }
	grid := make([][]cell, rows)
	for i := range grid {
		grid[i] = make([]cell, cols)
	}

	var buf strings.Builder
	buf.Grow(cols * rows * 6)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		newCols, newRows := termfx.GetSize()
		if newCols != cols || newRows != rows {
			cols, rows = newCols, newRows
			flakes = make([]snowflake, cols/3)
			for i := range flakes {
				flakes[i] = newFlake(rng, cols, rows, true)
			}
			grid = make([][]cell, rows)
			for i := range grid {
				grid[i] = make([]cell, cols)
			}
		}

		for i := range grid {
			for j := range grid[i] {
				grid[i][j].r = 0
			}
		}

		for i := range flakes {
			f := &flakes[i]
			f.y += f.speed
			f.x += f.drift
			if f.x < 0 {
				f.x = float64(cols - 1)
			} else if f.x >= float64(cols) {
				f.x = 0
			}
			if f.y >= float64(rows) {
				flakes[i] = newFlake(rng, cols, rows, false)
				continue
			}
			r, c := int(f.y), int(f.x)
			if r >= 0 && r < rows && c >= 0 && c < cols {
				grid[r][c].r = f.char
			}
		}

		buf.Reset()
		buf.WriteString("\033[H")
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				ch := grid[i][j].r
				if ch != 0 {
					buf.WriteString("\033[97m")
					buf.WriteRune(ch)
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
		time.Sleep(80 * time.Millisecond)
	}
}

func newFlake(rng *rand.Rand, cols, rows int, randomY bool) snowflake {
	y := 0.0
	if randomY {
		y = rng.Float64() * float64(rows)
	}
	return snowflake{
		x:     rng.Float64() * float64(cols),
		y:     y,
		speed: 0.3 + rng.Float64()*0.7,
		drift: (rng.Float64()*0.6 - 0.3),
		char:  snowChars[rng.Intn(len(snowChars))],
	}
}

func init() {
	termfx.Register("Snow", &snow{})
}

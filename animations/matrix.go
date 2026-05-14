package animations

import (
	"context"
	"github.com/mohamedation/gotermfx/termfx"
	"math/rand"
	"os"
	"strings"
	"time"
)

type matrix struct{}

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%&")

type colState struct {
	pos    int
	length int
	speed  int
	tick   int
}

func (m *matrix) Run(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cols, rows := termfx.GetSize()
	state := makeState(rng, cols, rows)

	var buf strings.Builder
	buf.Grow(cols * rows * 12)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		newCols, newRows := termfx.GetSize()
		if newCols != cols || newRows != rows {
			cols, rows = newCols, newRows
			state = makeState(rng, cols, rows)
			buf.Reset()
			buf.Grow(cols * rows * 12)
		}

		buf.Reset()
		buf.WriteString("\033[H")

		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				s := &state[j]
				dist := s.pos - i
				switch {
				case dist == 0 && s.pos >= 0:
					// bright white
					buf.WriteString("\033[97;1m")
					buf.WriteRune(charset[rng.Intn(len(charset))])
					buf.WriteString("\033[0m")
				case dist > 0 && dist <= s.length/3:
					buf.WriteString("\033[92m") // bright green
					buf.WriteRune(charset[rng.Intn(len(charset))])
					buf.WriteString("\033[0m")
				case dist > 0 && dist <= 2*s.length/3:
					buf.WriteString("\033[32m") // green
					buf.WriteRune(charset[rng.Intn(len(charset))])
					buf.WriteString("\033[0m")
				case dist > 0 && dist <= s.length:
					buf.WriteString("\033[2;32m") // dim green
					buf.WriteRune(charset[rng.Intn(len(charset))])
					buf.WriteString("\033[0m")
				default:
					buf.WriteByte(' ')
				}
			}
			if i < rows-1 {
				buf.WriteString("\r\n")
			}
		}

		os.Stdout.WriteString(buf.String())

		for j := range state {
			state[j].tick++
			if state[j].tick%state[j].speed == 0 {
				state[j].pos++
				if state[j].pos > rows+state[j].length {
					state[j].pos = -(4 + rng.Intn(rows/2))
					state[j].length = 4 + rng.Intn(rows/2+1)
					state[j].speed = 1 + rng.Intn(3)
				}
			}
		}

		time.Sleep(40 * time.Millisecond)
	}
}

func makeState(rng *rand.Rand, cols, rows int) []colState {
	s := make([]colState, cols)
	for j := range s {
		s[j] = colState{
			pos:    -(1 + rng.Intn(rows)),
			length: 4 + rng.Intn(rows/2+1),
			speed:  1 + rng.Intn(3),
			tick:   rng.Intn(100),
		}
	}
	return s
}

func init() {
	termfx.Register("Matrix", &matrix{})
}

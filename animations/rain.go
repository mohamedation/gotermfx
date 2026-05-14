package animations

import (
	"context"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mohamedation/gotermfx/termfx"
)

type rain struct{}

func (r *rain) Run(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var buf strings.Builder
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		cols, rows := termfx.GetSize()
		buf.Reset()
		buf.Grow(cols * rows * 5)
		buf.WriteString("\033[H")
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				if rng.Float32() < 0.02 {
					buf.WriteString("\033[36m|\033[0m")
				} else {
					buf.WriteByte(' ')
				}
			}
			if i < rows-1 {
				buf.WriteString("\r\n")
			}
		}
		os.Stdout.WriteString(buf.String())
		time.Sleep(60 * time.Millisecond)
	}
}

func init() {
	termfx.Register("Rain", &rain{})
}

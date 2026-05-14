package animations

import (
	"context"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mohamedation/gotermfx/termfx"
)

type fireworks struct{}

type rocket struct {
	x, y      float64
	targetY   float64
	exploded  bool
	particles []particle
	color     string
}

type particle struct {
	x, y   float64
	vx, vy float64
	life   float64
	char   rune
}

var rocketColors = []string{
	"\033[91m", // red
	"\033[93m", // yellow
	"\033[92m", // green
	"\033[96m", // cyan
	"\033[95m", // magenta
	"\033[97m", // white
	"\033[94m", // blue
}

var particleChars = []rune{'*', '+', '.', '·', '°', '•', '✦', '✧'}

func (fw *fireworks) Run(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cols, rows := termfx.GetSize()

	rockets := make([]*rocket, 0, 8)

	grid := make([]fwCell, cols*rows)

	var buf strings.Builder
	buf.Grow(cols * rows * 10)

	launchTick := 0

	onceLaunched := false

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		newCols, newRows := termfx.GetSize()
		if newCols != cols || newRows != rows {
			cols, rows = newCols, newRows
			rockets = rockets[:0]
			grid = make([]fwCell, cols*rows)
		}

		if termfx.IsOnce(ctx) && !onceLaunched {
			onceLaunched = true
			for i := 0; i < 3; i++ {
				rockets = append(rockets, &rocket{
					x:       float64(5 + rng.Intn(cols-10)),
					y:       float64(rows - 1),
					targetY: float64(3 + rng.Intn(rows/2)),
					color:   rocketColors[rng.Intn(len(rocketColors))],
				})
			}
		}

		launchTick++
		if !termfx.IsOnce(ctx) && launchTick >= 15+rng.Intn(15) && len(rockets) < 6 {
			launchTick = 0
			r := &rocket{
				x:       float64(5 + rng.Intn(cols-10)),
				y:       float64(rows - 1),
				targetY: float64(3 + rng.Intn(rows/2)),
				color:   rocketColors[rng.Intn(len(rocketColors))],
			}
			rockets = append(rockets, r)
		}

		for i := range grid {
			grid[i].r = 0
		}

		alive := rockets[:0]
		for _, r := range rockets {
			if !r.exploded {
				r.y -= 1.5
				if r.y <= r.targetY {
					r.exploded = true
					explode(rng, r, cols, rows)
				} else {
					plot(grid, cols, rows, int(r.x), int(r.y), '|', r.color)
					plot(grid, cols, rows, int(r.x), int(r.y)-1, '^', r.color)
				}
			}

			if r.exploded {
				anyAlive := false
				for i := range r.particles {
					p := &r.particles[i]
					if p.life <= 0 {
						continue
					}
					anyAlive = true
					p.x += p.vx
					p.y += p.vy
					p.vy += 0.05 // gravity
					p.vx *= 0.97 // drag
					p.life -= 0.04

					alpha := p.life
					var color string
					switch {
					case alpha > 0.6:
						color = r.color
					case alpha > 0.3:
						color = "\033[2m" + r.color // dim
					default:
						color = "\033[2;90m" // very dim gray
					}
					plot(grid, cols, rows, int(p.x), int(p.y), p.char, color)
				}
				if anyAlive {
					alive = append(alive, r)
				}
			} else {
				alive = append(alive, r)
			}
		}
		rockets = alive

		if termfx.IsOnce(ctx) && len(rockets) == 0 {
			return
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
		time.Sleep(50 * time.Millisecond)
	}
}

func explode(rng *rand.Rand, r *rocket, cols, rows int) {
	count := 30 + rng.Intn(30)
	r.particles = make([]particle, count)
	for i := range r.particles {
		angle := rng.Float64() * 2 * math.Pi
		speed := 0.5 + rng.Float64()*2.0
		r.particles[i] = particle{
			x:    r.x,
			y:    r.y,
			vx:   math.Cos(angle) * speed,
			vy:   math.Sin(angle) * speed * 0.5,
			life: 0.8 + rng.Float64()*0.2,
			char: particleChars[rng.Intn(len(particleChars))],
		}
	}
}

type fwCell struct {
	r     rune
	color string
}

func plot(grid []fwCell, cols, rows, x, y int, ch rune, color string) {
	if x < 0 || x >= cols || y < 0 || y >= rows {
		return
	}
	grid[y*cols+x].r = ch
	grid[y*cols+x].color = color
}

func init() {
	termfx.Register("Fireworks", &fireworks{})
}

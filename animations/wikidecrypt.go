package animations

import (
	"context"
	"encoding/json"
	"github.com/mohamedation/gotermfx/termfx"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type wikiDecrypt struct{}

type wikiSummary struct {
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

type wikiCell struct {
	r      rune
	locked bool
	header bool
}

var wikiScrambleChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*+-=~<>?/|")

func fetchRandomWiki(ctx context.Context) (*wikiSummary, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://en.wikipedia.org/api/rest_v1/page/random/summary", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "github.com/mohamedation/gotermfx/1.0 terminal-animation")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var s wikiSummary
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func wrapWords(text string, maxWidth int) []string {
	if maxWidth < 1 {
		maxWidth = 80
	}
	var out []string
	for _, para := range strings.Split(text, "\n") {
		words := strings.Fields(para)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}
		cur := ""
		for _, w := range words {
			if cur == "" {
				cur = w
			} else if len(cur)+1+len(w) <= maxWidth {
				cur += " " + w
			} else {
				out = append(out, cur)
				cur = w
			}
		}
		if cur != "" {
			out = append(out, cur)
		}
	}
	return out
}

// todo> fix the connection animation
func wikiConnect(ctx context.Context) (*wikiSummary, error) {
	fetchResult := make(chan *wikiSummary, 1)
	fetchErr := make(chan error, 1)
	go func() {
		s, err := fetchRandomWiki(ctx)
		if err != nil {
			fetchErr <- err
		} else {
			fetchResult <- s
		}
	}()

	messages := []string{
		"> INITIALISING CONNECTION",
		"> LOCATING ENDPOINT",
		"> PERFORMING TLS HANDSHAKE",
		"> SENDING REQUEST",
		"> AWAITING RESPONSE",
		"> PARSING DATA STREAM",
		"> DECODING PAYLOAD",
		"> ARTICLE RECEIVED",
	}

	var buf strings.Builder
	msgIdx := 0
	nextMsg := time.Now().Add(200 * time.Millisecond)
	log := make([]string, 0, len(messages))

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case s := <-fetchResult:
			return s, nil
		case err := <-fetchErr:
			return nil, err
		case <-ticker.C:
			if time.Now().After(nextMsg) && msgIdx < len(messages) {
				msg := messages[msgIdx]
				dots := strings.Repeat(".", 42-len(msg))
				if msgIdx < len(messages)-1 {
					log = append(log, "\033[32m"+msg+" "+dots+" OK\033[0m")
				} else {
					log = append(log, "\033[97;1m"+msg+" "+dots+" >>>\033[0m")
				}
				msgIdx++
				nextMsg = time.Now().Add(300 * time.Millisecond)
			}
			cols, rows := termfx.GetSize()
			buf.Reset()
			buf.WriteString("\033[2J\033[H")
			top := (rows - len(log) - 3) / 2
			for i := 0; i < top; i++ {
				buf.WriteString("\r\n")
			}
			title := "[ ACCESSING GLOBAL KNOWLEDGE NETWORK ]"
			pad := (cols - len(title)) / 2
			if pad < 0 {
				pad = 0
			}
			buf.WriteString(strings.Repeat(" ", pad))
			buf.WriteString("\033[97;1m")
			buf.WriteString(title)
			buf.WriteString("\033[0m\r\n\r\n")
			for _, line := range log {
				buf.WriteString("  ")
				buf.WriteString(line)
				buf.WriteString("\r\n")
			}
			os.Stdout.WriteString(buf.String())
		}
	}
}

func (w *wikiDecrypt) Run(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		summary, err := wikiConnect(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
				continue
			}
		}

		cols, rows := termfx.GetSize()
		maxCols := cols - 4
		header := "[ " + summary.Title + " ]"
		displayLines := []string{header, ""}
		displayLines = append(displayLines, wrapWords(summary.Extract, maxCols)...)
		if len(displayLines) > rows-1 {
			displayLines = displayLines[:rows-1]
		}

		cells := make([]wikiCell, 0, 512)
		lineEnds := make([]int, len(displayLines))
		for li, line := range displayLines {
			for _, r := range line {
				cells = append(cells, wikiCell{
					r:      r,
					locked: r == ' ' || r == '\t',
					header: li == 0,
				})
			}
			lineEnds[li] = len(cells)
		}

		deadline := time.Now().Add(1400 * time.Millisecond)
		for time.Now().Before(deadline) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			wikiRender(cells, displayLines, lineEnds, rng, cols)
			time.Sleep(40 * time.Millisecond)
		}

		order := make([]int, 0, len(cells))
		for i, c := range cells {
			if !c.locked {
				order = append(order, i)
			}
		}
		rng.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

		n := max(len(order), 1)
		delay := time.Duration(3500/n) * time.Millisecond
		delay = max(delay, 6*time.Millisecond)
		delay = min(delay, 40*time.Millisecond)

		for _, idx := range order {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cells[idx].locked = true
			wikiRender(cells, displayLines, lineEnds, rng, cols)
			time.Sleep(delay)
		}

		wordCount := len(strings.Fields(summary.Extract))
		readDur := time.Duration(float64(wordCount)/228.0*60.0) * time.Second
		if readDur < 5*time.Second {
			readDur = 5 * time.Second
		}
		if readDur > 60*time.Second {
			readDur = 60 * time.Second
		}
		wikiRender(cells, displayLines, lineEnds, rng, cols)
		select {
		case <-ctx.Done():
			return
		case <-time.After(readDur):
		}

		if termfx.IsOnce(ctx) {
			return
		}
	}
}

func wikiRender(cells []wikiCell, lines []string, lineEnds []int, rng *rand.Rand, cols int) {
	var buf strings.Builder
	buf.Grow(len(cells) * 10)
	buf.WriteString("\033[H")
	ci := 0
	for li, line := range lines {
		end := lineEnds[li]
		buf.WriteString("  ")
		for ci < end {
			c := cells[ci]
			if c.locked {
				if c.header {
					buf.WriteString("\033[97;1m") // bold white
				} else {
					buf.WriteString("\033[92m") // bright green
				}
				buf.WriteRune(c.r)
				buf.WriteString("\033[0m")
			} else {
				if rng.Intn(3) == 0 {
					buf.WriteString("\033[32m")
				} else {
					buf.WriteString("\033[2;32m")
				}
				buf.WriteRune(wikiScrambleChars[rng.Intn(len(wikiScrambleChars))])
				buf.WriteString("\033[0m")
			}
			ci++
		}
		written := len([]rune(line)) + 2
		for p := written; p < cols; p++ {
			buf.WriteByte(' ')
		}
		buf.WriteString("\r\n")
	}
	os.Stdout.WriteString(buf.String())
}

func init() {
	termfx.Register("WikiDecrypt", &wikiDecrypt{})
}

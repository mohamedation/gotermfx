package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/mohamedation/gotermfx/animations"
	"github.com/mohamedation/gotermfx/termfx"
)

func main() {
	helpFlag := flag.Bool("h", false, "show help")
	interactive := flag.Bool("i", false, "interactive mode")
	once := flag.Bool("1", false, "run the animation once then exit")
	duration := flag.Duration("d", 10*time.Second, "max run time in once mode for continuous animations (e.g. 10s, 1m)")
	flag.Usage = func() { printHelp() }
	flag.Parse()

	// we need to check duration for animations that need to finish a sequence
	durationSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "d" {
			durationSet = true
		}
	})
	args := flag.Args()

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	idx := -1
	if *interactive {
		fmt.Println("gotermfx Terminal Animation CLI")
		fmt.Println("Available animations:")
		for i, name := range termfx.List() {
			fmt.Printf("%d. %s\n", i+1, name)
		}
		fmt.Print("Select animation: ")
		var choice string
		fmt.Scanln(&choice)
		idx = resolveAnimation(strings.TrimSpace(choice))
	} else if len(args) == 0 {
		// random mode
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		idx = rng.Intn(len(termfx.List()))
	} else {
		idx = resolveAnimation(args[0])
	}

	if idx == -1 {
		fmt.Fprintf(os.Stderr, "error: unknown animation %q\n\n", args[0])
		printHelp()
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if *once {
		if durationSet {
			var timeoutCancel context.CancelFunc
			ctx, timeoutCancel = context.WithTimeout(ctx, *duration)
			defer timeoutCancel()
		}
		ctx = termfx.WithOnce(ctx)
	}

	// raw mode
	oldTermios, err := setRawMode()
	if err == nil {
		defer restoreMode(oldTermios)
	}

	// clear screen and no cursor
	os.Stdout.WriteString("\033[?25l\033[2J\033[H")
	defer os.Stdout.WriteString("\033[?25h\033[2J\033[H")

	// cancel on any keypress.
	go func() {
		buf := make([]byte, 1)
		os.Stdin.Read(buf)
		cancel()
	}()

	// cancel on SIGTERM, but since we have cancel on any keypress, this is mostly uneeded but i will keep it for now in case we need  change the any key method
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	termfx.Get(idx).Run(ctx)
}

func printHelp() {
	fmt.Println("gotermfx — terminal animations")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  gotermfx                    run a random animation")
	fmt.Println("  gotermfx <name|number>      run a specific animation")
	fmt.Println("  gotermfx -i                 interactive picker")
	fmt.Println("  gotermfx -1              run once then exit")
	fmt.Println("  gotermfx -1 -d 5s        set duration for continuous animations (default 10s)")
	fmt.Println("  gotermfx -h                 show this help")
	fmt.Println()
	fmt.Println("ANIMATIONS:")
	for i, name := range termfx.List() {
		fmt.Printf("  %-3d %s\n", i+1, name)
	}
	fmt.Println()
	fmt.Println("Press any key to stop an animation.")
}

func resolveAnimation(choice string) int {
	for i, name := range termfx.List() {
		if fmt.Sprintf("%d", i+1) == choice || strings.EqualFold(name, choice) {
			return i
		}
	}
	return -1
}

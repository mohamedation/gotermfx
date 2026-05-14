# gotermfx

A modular, zero-dependency terminal animation CLI written in Go.

## Animations

| Name       | Description                                         |
|------------|-----------------------------------------------------|
| `rain`     | Cyan ASCII rain drops falling down the screen       |
| `matrix`   | Green Matrix-style scrolling characters with trails |
| `snow`     | White snowflakes drifting and falling               |
| `starfield`| Warp-speed star tunnel zooming toward the viewer    |
| `fireworks`| Colourful rockets launching and bursting            |

## Requirements

- Go 1.21+
- macOS or Linux (uses POSIX `ioctl` for terminal size and raw mode)

## Build & Install

```sh
# Build binary into ./bin/gotermfx
make build

# Install to /usr/local/bin (may need sudo)
make install

# Remove installed binary
make uninstall

# Clean build artefacts
make clean
```

## Usage

```sh
# Run an animation by name (case-insensitive)
gotermfx matrix
gotermfx rain
gotermfx snow
gotermfx starfield
gotermfx fireworks

# Or by number
gotermfx 1

# Interactive picker
gotermfx -i
```

Press **any key** or **Ctrl-C** to exit.

## Adding a New Animation

1. Create `animations/myanim.go` in the `animations` package.
2. Define a struct and implement `Run(ctx context.Context)`.
3. Call `Register("Name", &myAnim{})` inside an `init()` function.
4. That's it — it will appear in the list automatically.

```go
package animations

import (
    "context"
    "os"
    "time"
)

type myAnim struct{}

func (a *myAnim) Run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
        }
        cols, rows := GetSize()
        _ = cols
        _ = rows
        os.Stdout.WriteString("\033[H" + "Hello!\r\n")
        time.Sleep(100 * time.Millisecond)
    }
}

func init() {
    Register("MyAnim", &myAnim{})
}
```

## Project Structure

```
gotermfx/
├── main.go              # CLI entry point, flag parsing, raw mode setup
├── term_darwin.go       # macOS raw terminal mode (ioctl TIOCGETA/TIOCSETA)
├── term_linux.go        # Linux raw terminal mode (ioctl TCGETS/TCSETS)
├── animations/
│   ├── animations.go    # Animation interface + registry
│   ├── term.go          # GetSize() via ioctl TIOCGWINSZ (unix build tag)
│   ├── rain.go
│   ├── matrix.go
│   ├── snow.go
│   ├── starfield.go
│   └── fireworks.go
├── go.mod
├── Makefile
└── README.md
```

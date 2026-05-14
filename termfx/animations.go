package termfx

import "context"

// Animation is the interface every animation must implement
// run should block until ctx is cancelled
type Animation interface {
	Run(ctx context.Context)
}

// once flag handling
type onceKey struct{}

func WithOnce(ctx context.Context) context.Context {
	return context.WithValue(ctx, onceKey{}, true)
}

func IsOnce(ctx context.Context) bool {
	v, _ := ctx.Value(onceKey{}).(bool)
	return v
}

var registry []Animation
var names []string

func Register(name string, a Animation) {
	names = append(names, name)
	registry = append(registry, a)
}

func List() []string {
	return names
}

func Get(i int) Animation {
	return registry[i]
}

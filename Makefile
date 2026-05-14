BINARY   := gotermfx
BINDIR   := bin
INSTALLDIR := /usr/local/bin

GO       := go
GOFLAGS  := -trimpath -ldflags="-s -w"

.PHONY: all build run install uninstall clean

all: build

## build: compile the binary into ./bin/gotermfx
build:
	@mkdir -p $(BINDIR)
	$(GO) build $(GOFLAGS) -o $(BINDIR)/$(BINARY) .
	@echo "Built: $(BINDIR)/$(BINARY)"

## run: build and run with ARGS (e.g. make run ARGS=matrix)
run: build
	$(BINDIR)/$(BINARY) $(ARGS)

## install: install the binary to /usr/local/bin
install: build
	install -m 0755 $(BINDIR)/$(BINARY) $(INSTALLDIR)/$(BINARY)
	@echo "Installed: $(INSTALLDIR)/$(BINARY)"

## uninstall: remove the installed binary
uninstall:
	rm -f $(INSTALLDIR)/$(BINARY)
	@echo "Removed: $(INSTALLDIR)/$(BINARY)"

## clean: remove build artefacts
clean:
	rm -rf $(BINDIR)
	@echo "Cleaned."

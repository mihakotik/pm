PREFIX ?= /usr
DESTDIR ?=
BINDIR ?= $(PREFIX)/bin

GO ?= go
GOFLAGS ?= -ldflags="-s -w"
BIN := bin/pm

.PHONY: all build install clean test fmt vet

all: build

build:
	$(GO) build $(GOFLAGS) -o $(BIN) ./cmd/pm

install: build
	install -Dm755 $(BIN) $(DESTDIR)$(BINDIR)/pm

clean:
	rm -rf bin/

test:
	$(GO) test -v ./...

fmt:
	gofmt -s -w .

vet:
	$(GO) vet ./...

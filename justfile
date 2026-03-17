
# this list
default:
  just --list

# build the binary
build:
  go build -ldflags "-s -w -X main.version=dev -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o propel

# run it
run: build
  ./propel

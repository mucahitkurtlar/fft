# fft

_find font types_

Find font types used on a website.

## Installation

Install latest version from [releases](https://github.com/mucahitkurtlar/fft/releases) or build from source.

## Build from source

### Pull source code

```bash
git clone https://github.com/mucahitkurtlar/fft.git
cd fft
go mod download
go build
```

### Usage

```
fft <url> [flags]

Examples:
        fft https://www.example.com
        fft https://www.example.com -m 20

Flags:
  -g, --go-routines uint8        number of go routines to use (default 3)
  -t, --go-to-timeout float      timeout for page navigation (ms) (default 30000)
  -h, --help                     help for fft
  -m, --max-pages uint16         maximum number of pages to crawl (default 1)
  -n, --net-idle-timeout float   timeout for network idle (ms) (default 30000)

```

### Example
```bash
./fft https://www.example.com
```

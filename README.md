# brotherql

Go library and CLI for printing to **Brother QL** label printers over USB.

> **Status**: Early development (v0.x). API may change before v1.0.

## Acknowledgements

Inspired by the [`brother_ql`](https://github.com/pklaus/brother_ql) Python library by Philipp Klaus, which served as a protocol reference and inspiration for this Go port.

## Features

- Idiomatic Go API: pass an `image.Image`, get a printed label.
- CLI tool for ad-hoc printing and shell-script integration.
- Native USB communication via [gousb](https://github.com/google/gousb) — no CUPS or system driver required.
- Status query (paper present, errors, media dimensions).
- Tested via golden files cross-validated against the Python reference.

## Supported hardware

- Brother QL-700 over USB.
- Brother QL-710W over USB (wireless interface not used; USB only).

Other Brother QL models (QL-720NW, QL-800, QL-1100, QL-820NWB, etc.) are not yet supported in v1. Contributions welcome.

## Installation

### Prerequisites

`libusb` and `pkg-config` must be installed on your system:

| OS | Install |
|----|---------|
| macOS | `brew install libusb pkg-config` |
| Debian/Ubuntu | `sudo apt install libusb-1.0-0-dev pkg-config` |
| Fedora | `sudo dnf install libusb1-devel pkgconf-pkg-config` |

### As a library

```bash
go get github.com/kurtinge/brotherql
```

### As a CLI

```bash
go install github.com/kurtinge/brotherql/cmd/brotherql@latest
```

Pre-built binaries are also available on the [Releases](https://github.com/kurtinge/brotherql/releases) page.

## Quick start (library)

```go
package main

import (
    "image/png"
    "log"
    "os"

    "github.com/kurtinge/brotherql"
)

func main() {
    p, err := brotherql.Open()
    if err != nil {
        log.Fatal(err)
    }
    defer p.Close()

    f, err := os.Open("label.png")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    img, err := png.Decode(f)
    if err != nil {
        log.Fatal(err)
    }

    err = p.Print(img, brotherql.PrintOptions{
        Label:   brotherql.Label62,
        AutoCut: true,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Quick start (CLI)

List connected printers:

```bash
brotherql list
# QL-700 serial=000D5G103245 (bus 2 addr 5)
```

Show printer status:

```bash
brotherql status
# Ready: yes
# Media: 62mm continuous
# Error: none
```

Print a label:

```bash
brotherql print --label 62 image.png
```

Print to a specific printer (multiple connected):

```bash
brotherql print --label 62 --serial 000D5G103245 image.png
```

## Supported labels

| Constant | Description |
|----------|-------------|
| `Label62` | 62mm continuous tape (DK-22205 etc.) |
| `Label62x29` | 62mm × 29mm die-cut |

More to come; PRs welcome.

## Testing

```bash
# Pure unit tests (no hardware)
go test ./...

# Hardware integration tests (requires connected printer)
go test -tags=hardware ./...
```

The test suite uses golden files cross-validated against the Python `brother_ql` reference. To regenerate them, see [`testdata/bootstrap.sh`](testdata/bootstrap.sh).

## Contributing

Open an issue first for non-trivial changes. PRs welcome — please include tests.

## License

MIT — see [LICENSE](LICENSE).

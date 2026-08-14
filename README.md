# goqrcode - Go QR Code Generator Library, CLI Tool and Web Service

[![pschlump](https://circleci.com/gh/pschlump/goqrcode.svg?style=svg)](https://app.circleci.com/pipelines/github/pschlump/goqrcode)

This is a fork of [skip2/go-qrcode](https://github.com/skip2/go-qrcode)
(originally by Tom Harwood) that has been extended with:

- **SVG output** — via the companion library
  [goqrsvg](https://github.com/pschlump/goqrsvg), both in the CLI and the web
  service.
- **Text art output** — render a scannable QR Code directly in the terminal
  (`ToString` / `ToSmallString`).
- **A QR Code web service** — `qr-gen-server`, an HTTP server that generates
  PNG, SVG or text QR Codes from URL parameters.
- **A passphrase encryption helper** — the `crypt` package, AES-256-GCM
  encryption used to embed encrypted data in QR Codes.

Package `goqrcode` implements a QR Code encoder. A QR Code is a matrix
(two-dimensional) barcode. Arbitrary content may be encoded, with URLs being a
popular choice.

Each QR Code contains error recovery information to aid reading damaged or
obscured codes. There are four levels of error recovery: `Low`, `Medium`,
`High` and `Highest`. QR Codes with a higher recovery level are more robust to
damage, at the cost of being physically larger.

## Repository layout

| Path             | Contents                                                        |
|------------------|-----------------------------------------------------------------|
| `.`              | The `goqrcode` library package                                  |
| `qrcode/`        | Command line tool `qrcode`                                      |
| `qr-gen-server/` | HTTP QR Code generation service                                 |
| `crypt/`         | AES-256-GCM passphrase encryption helpers                       |
| `bitset/`        | Bit array used by the encoder                                   |
| `reedsolomon/`   | Reed-Solomon error correction                                   |

## Install

    go get github.com/pschlump/goqrcode/...

## Library usage

    import goqrcode "github.com/pschlump/goqrcode"

- **Create a PNG image:**

        var png []byte
        png, err := goqrcode.Encode("https://example.org", goqrcode.Medium, 256)

- **Create a PNG image and write to a file:**

        err := goqrcode.WriteFile("https://example.org", goqrcode.Medium, 256, "qr.png")

- **Create a PNG image with custom colors and write to file:**

        err := goqrcode.WriteColorFile("https://example.org", goqrcode.Medium, 256, color.Black, color.White, "qr.png")

- **Print a QR Code as terminal text art:**

        q, err := goqrcode.New("https://example.org", goqrcode.Medium)
        fmt.Println(q.ToString(false))

All examples use the `goqrcode.Medium` error recovery level and create a fixed
256x256px QR Code. The `WriteColorFile` example creates a white-on-black
instead of black-on-white QR Code.

To generate a variable sized image instead, specify a negative size (in place
of the 256 above), such as -5. A size of -5 sets each module (QR Code "pixel")
to be 5px wide/high, producing the minimum size image with no extra padding.

The maximum capacity of a QR Code varies according to the content encoded and
the error recovery level. The maximum capacity is 2,953 bytes, 4,296
alphanumeric characters, 7,089 numeric digits, or a combination of these.

More runnable examples are in [example_test.go](example_test.go); run them
with `go test -run Example .`.

## Command line tool

    cd ./qrcode
    go build
    ./qrcode -o test1.png "otpauth://totp/www.2c-why.com:pschlump@gmail.com?secret=4S62BZNFXXSZLCRO&issuer=www.2c-why.com"

```
qrcode -- QR Code encoder in Go
https://github.com/pschlump/goqrcode

Flags:
  -i        invert black and white
  -l string
        Level of error redundancy, h|m|l for high, medium, low (default "h")
  -o string
        out PNG file prefix, empty for stdout
  -s int
        image size (pixel) (default 256)
  -svg
        print as SVG on stdout
  -txt
        print as text art on stdout
```

Examples:

    # PNG to a file (format is chosen by file extension or flag)
    ./qrcode -o id_1234467721.png "http://www.mysite.com/qr-display/1234467721"

    # SVG to a file
    ./qrcode -o xyz.svg -svg "http://www.mysite.com/qr-display/1234467721"

    # Text art to the terminal
    ./qrcode -txt "http://www.mysite.com/qr-display/1234467721"

## Web service: qr-gen-server

`qr-gen-server` serves generated QR Codes over HTTP. Build and run it:

    cd ./qr-gen-server
    go build
    ./qr-gen-server --hostport :2003 --cfg cfg.json

Then request QR Codes by URL:

| Parameter     | Values                   | Default | Description                          |
|---------------|--------------------------|---------|--------------------------------------|
| `url`         | any text                 | (none)  | Content to encode                    |
| `fmt`         | `png`, `svg`, `text`     | `png`   | Output format                        |
| `redundancy`  | `H`, `M`, `L`            | `M`     | Error recovery level                 |
| `size`        | integer                  | `256`   | Image size in pixels (PNG only)      |
| `invert`      | boolean                  | `false` | Swap foreground/background colors     |

Examples:

    # PNG (default)
    wget -O qr.png 'http://127.0.0.1:2003/?url=http://www.mysite.com/xyz'

    # SVG
    wget -O qr.svg 'http://127.0.0.1:2003/?url=http://www.mysite.com/xyz&fmt=svg'

    # Inverted colors, high redundancy text art
    wget -O - 'http://127.0.0.1:2003/?url=hello&fmt=text&invert=true&redundancy=H'

## crypt package

The `crypt` package encrypts data with a passphrase (AES-256-GCM, key derived
from an MD5 digest of the passphrase), so that encrypted content can be
embedded in a QR Code:

    ciphertext, err := crypt.EncryptWithPassphrase([]byte("secret"), "passphrase")
    plaintext, err := crypt.DecryptWithPassphrase(ciphertext, "passphrase")

## Development

    make test    # go vet + go test ./...
    make lint    # golangci-lint run ./...
    make all    # build everything

## Links

- [http://en.wikipedia.org/wiki/QR_code](http://en.wikipedia.org/wiki/QR_code)
- [ISO/IEC 18004:2006](http://www.iso.org/iso/catalogue_detail.htm?csnumber=43655) - Main QR Code specification
- [https://github.com/qpliu/qrencode-go/](https://github.com/qpliu/qrencode-go/) - alternative Go QR encoding library based on [ZXing](https://github.com/zxing/zxing)
- [goqrsvg](https://github.com/pschlump/goqrsvg) - SVG rendering companion library

# goqrcode

This project was cloned because the original author seems to have abandoned it.

A fixed version based on "http://github.com/pschlump/qrcode".  This project did most of the
original work.

Changed command line to work slightly differently.

Implemented .SVG output.

Implemented .JPG output.

Package `goqrcode` implements a QR Code encoder. 

A QR Code is a matrix (two-dimensional) barcode. Arbitrary content may be encoded, with URLs being a popular choice :)

Each QR Code contains error recovery information to aid reading damaged or obscured codes. There are four levels of error recovery: Low, medium, high and highest. QR Codes with a higher recovery level are more robust to damage, at the cost of being physically larger.

## Install

    go get -u github.com/pschlump/qrcode/...

A command-line tool `qrcode` will be built into `$GOPATH/bin/`.

## Usage

    import qrcode "github.com/pschlump/qrcode"

- **Create a PNG image:**

        var png []byte
        png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)

- **Create a PNG image and write to a file:**

        err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "qr.png")

- **Create a PNG image with custom colors and write to file:**

        err := qrcode.WriteColorFile("https://example.org", qrcode.Medium, 256, color.Black, color.White, "qr.png")

All examples use the qrcode.Medium error Recovery Level and create a fixed
256x256px size QR Code. The last function creates a white on black instead of black
on white QR Code.

The maximum capacity of a QR Code varies according to the content encoded and
the error recovery level. The maximum capacity is 2,953 bytes, 4,296
alphanumeric characters, 7,089 numeric digits, or a combination of these.

## Links

- [http://en.wikipedia.org/wiki/QR_code](http://en.wikipedia.org/wiki/QR_code)
- [ISO/IEC 18004:2006](http://www.iso.org/iso/catalogue_detail.htm?csnumber=43655) - Main QR Code specification (approx CHF 198,00)<br>
- [https://github.com/qpliu/qrencode-go/](https://github.com/qpliu/qrencode-go/) - alternative Go QR encoding library based on [ZXing](https://github.com/zxing/zxing)

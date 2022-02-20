
all:
	go build

gen_code:
	mkdir -p ./params
	/Users/philip/go/src/github.com/pschlump/req_param/req_param -i gen_qr.json -p ApiQRGen >params/gen_qr.go
	/Users/philip/go/bin/goimports -w params/gen_qr.go


all:
	go build

# REQ_PARAM=/Users/philip/go/src/github.com/pschlump/req_param/req_param
REQ_PARAM=~/bin/req_param

gen_code:
	mkdir -p ./params
	${REQ_PARAM} -i gen_qr.json -p ApiQRGen >params/gen_qr.go
	/Users/philip/go/bin/goimports -w params/gen_qr.go

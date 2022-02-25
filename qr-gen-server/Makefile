
all:
	go build

# REQ_PARAM=/Users/philip/go/src/github.com/pschlump/req_param/req_param
REQ_PARAM=~/bin/req_param

gen_code:
	mkdir -p ./params
	${REQ_PARAM} -i gen_qr.json -p ApiQRGen >params/gen_qr.go
	/Users/philip/go/bin/goimports -w params/gen_qr.go

test: test1 test2
	@echo PASS

test1:
	mkdir -p out
	wget -o out/abc.def.err -O out/abc.def.svg 'http://127.0.0.1:2003/?url=abc.def&fmt=svg'
	diff out/abc.def.svg ref

test2:
	mkdir -p out
	wget -o out/abc.xyz.err -O out/abc.xyz.png 'http://127.0.0.1:2003/?url=abc.xyz'
	diff out/abc.xyz.png ref

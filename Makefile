
all:
	go build


test001:
	go test
	diff ./out/example.png ./ref/example.png



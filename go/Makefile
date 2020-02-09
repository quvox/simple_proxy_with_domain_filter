
build:
	go build -o main *.go

build-debug:
	go build -tags=debug -o main *.go


clean:
	rm -f main


run: build
	./main -block ../conf/blocklist.txt



.PHONY: clean
clean:
	rm -f ./gol

.PHONY: build
build: clean
	go build -o ./gol main.go

.PHONY: install
install: build
	chmod u+x ./gol
	mv ./gol /usr/local/bin/gol


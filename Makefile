run:
	go run .

build:
	mkdir -p builds
	go build -o builds/koko

clean:
	rm -f builds/koko

.PHONY: run build clean

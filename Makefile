run: clean build-static build-sass run-go

deploy: build push

build: clean build-static build-sass build-go

clean:
	rm -fR build

run-go:
	go run .

build-go:
	go run . -build

build-sass:
	mkdir -p build/stylesheets
	sassc source/stylesheets/all.scss > build/stylesheets/all.css

build-static:
	mkdir -p build/font
	mkdir -p build/javascript
	cp source/favicon.png build/favicon.png
	cp -r source/font/* build/font/
	cp -r source/javascript/* build/javascript/

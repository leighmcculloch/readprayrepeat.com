run: clean build-static build-sass run-go

deploy: build push

build: clean build-static build-sass build-go

clean:
	rm -fR build

run-go:
	go run .

build-go:
	go run . -build

build-dirs:
	mkdir -p build/stylesheets
	mkdir -p build/font
	mkdir -p build/javascript

build-sass: build-dirs
	sassc source/stylesheets/all.scss > build/stylesheets/all.css

build-static: build-dirs
	cp source/favicon.png build/favicon.png
	cp -r source/font/* build/font/
	cp -r source/javascript/* build/javascript/

push:
	firebase login --no-localhost
	firebase deploy

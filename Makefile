run: clean static sass run-go

debug: clean static sass debug-go

build: clean static sass build-go

deploy: build push cdn

clean:
	rm -fR build
	rm -fR build-assets
	rm -f app

build-go:
	go build -o app
	./app -build

run-go:
	go build -o app
	./app

debug-go:
	godebug run *.go

sass:
	mkdir -p build-assets/stylesheets
	sassc source/stylesheets/all.scss > build-assets/stylesheets/all.css

static:
	mkdir -p build-assets/font
	mkdir -p build-assets/javascript
	cp source/favicon.png build-assets/favicon.png
	cp -r source/font/* build-assets/font/
	cp -r source/javascript/* build-assets/javascript/

push:
	cd build && \
		gsutil -m -h "Content-Type: text/html" cp -a public-read -r . gs://today.bible/ 
	gsutil -q -m -h "Content-Type: text/css"                cp -Z -a public-read -r `find build-assets/stylesheets -type f | grep    '\.css'`                     gs://today.bible/stylesheets/
	gsutil -q -m -h "Content-Type: application/javascript"  cp -Z -a public-read -r `find build-assets/javascript  -type f | grep    '\.js'`                      gs://today.bible/javascript/
	gsutil -q -m -h "Content-Type:"                         cp -Z -a public-read -r `find build-assets/favicon.png -type f `                                      gs://today.bible/
	gsutil -q -m -h "Content-Type:"                         cp -Z -a public-read -r `find build-assets/font        -type f | grep    '\.\(svg\|eot\|png\)'`       gs://today.bible/font/
	gsutil -q -m -h "Content-Type: application/x-font-woff" cp -Z -a public-read -r `find build-assets/font        -type f | grep    '\.woff'`                    gs://today.bible/font/
	gsutil -q -m -h "Content-Type: font/truetype"           cp -Z -a public-read -r `find build-assets/font        -type f | grep    '\.ttf'`                     gs://today.bible/font/
	gsutil web set -m index.html -e 404.html gs://today.bible

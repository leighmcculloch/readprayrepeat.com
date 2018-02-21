export CLOUDFLARE_ZONE = 462201475c5da00f8d7a3bd778226fe5

run: clean static sass run-go

debug: clean static sass debug-go

build: clean static sass build-go

deploy: build push cdn

clean:
	rm -fR build
	rm -fR build-assets
	rm -f app

build-go:
	vgo build -o app
	./app -build

run-go:
	vgo build -o app
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
	gsutil -q -m -h "Content-Type:"                         cp -Z -a public-read -r `find build-assets/font        -type f | grep    '\.\(svg\|eot\|png\)'`       gs://today.bible/font/
	gsutil -q -m -h "Content-Type: application/x-font-woff" cp -Z -a public-read -r `find build-assets/font        -type f | grep    '\.woff'`                    gs://today.bible/font/
	gsutil -q -m -h "Content-Type: font/truetype"           cp -Z -a public-read -r `find build-assets/font        -type f | grep    '\.ttf'`                     gs://today.bible/font/
	gsutil web set -m index.html -e 404.html gs://today.bible

cdn:
	curl -X DELETE "https://api.cloudflare.com/client/v4/zones/$(CLOUDFLARE_ZONE)/purge_cache" \
		-H "X-Auth-Email: $(CLOUDFLARE_EMAIL)" \
		-H "X-Auth-Key: $(CLOUDFLARE_CLIENT_API_KEY)" \
		-H "Content-Type: application/json" \
		--data '{"purge_everything":true}'

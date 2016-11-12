export CLOUDFLARE_ZONE = 68fa298d917e2222f13f0afe848917ba

run: clean static sass run-go

debug: clean static sass debug-go

build: clean static sass build-go

deploy: build push cdn

clean:
	rm -fR build
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
	wt compile -b build/stylesheets/ source/stylesheets/all.scss

static:
	mkdir -p build/font
	mkdir -p build/javascript
	cp source/CNAME build/CNAME
	cp source/favicon.png build/favicon.png
	cp -r source/font/* build/font/
	cp -r source/javascript/* build/javascript/

push:
	gsutil -q -m -h "Content-Type: text/html"               cp -Z -a public-read -r `find build             -type f | grep -v '\.'`                        gs://readprayrepeat.com/ 
	gsutil -q -m -h "Content-Type: text/html"               cp -Z -a public-read -r `find build             -type f | grep    '\.\(html\|net\|esv\|gnt\)'` gs://readprayrepeat.com/ 
	gsutil -q -m -h "Content-Type: text/css"                cp -Z -a public-read -r `find build/stylesheets -type f | grep    '\.css'`                     gs://readprayrepeat.com/stylesheets/
	gsutil -q -m -h "Content-Type: application/javascript"  cp -Z -a public-read -r `find build/javascript  -type f | grep    '\.js'`                      gs://readprayrepeat.com/javascript/
	gsutil -q -m -h "Content-Type:"                         cp -Z -a public-read -r `find build/font        -type f | grep    '\.\(svg\|eot\|png\)'`       gs://readprayrepeat.com/font/
	gsutil -q -m -h "Content-Type: application/x-font-woff" cp -Z -a public-read -r `find build/font        -type f | grep    '\.woff'`                    gs://readprayrepeat.com/font/
	gsutil -q -m -h "Content-Type: font/truetype"           cp -Z -a public-read -r `find build/font        -type f | grep    '\.ttf'`                     gs://readprayrepeat.com/font/
	gsutil web set -m index.html -e 404.html gs://readprayrepeat.com

cdn:
	curl -X DELETE "https://api.cloudflare.com/client/v4/zones/$(CLOUDFLARE_ZONE)/purge_cache" \
		-H "X-Auth-Email: $(CLOUDFLARE_EMAIL)" \
		-H "X-Auth-Key: $(CLOUDFLARE_CLIENT_API_KEY)" \
		-H "Content-Type: application/json" \
		--data '{"purge_everything":true}'

setup:
	go get github.com/wellington/wellington/wt


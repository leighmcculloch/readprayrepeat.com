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
	mkdir -p build/stylesheets
	cd source; wt compile -b ../build stylesheets/all.scss

static:
	mkdir -p build/font
	mkdir -p build/javascript
	cp source/CNAME build/CNAME
	cp source/favicon.png build/favicon.png
	cp -r source/font/* build/font/
	cp -r source/javascript/* build/javascript/

push: push-s3

push-s3:
	aws s3 sync build/ s3://readprayrepeat.com/ --quiet --acl public-read --exclude "*" --include "*.eot" --include "*.svg" --include "*.ttf" --include "*.woff" --include "*.js" --include "*.css" --include "*.png"
	aws s3 sync build/ s3://readprayrepeat.com/ --quiet --acl public-read --content-type text/html --exclude "*.eot" --exclude "*.svg" --exclude "*.ttf" --exclude "*.woff" --exclude "*.js" --exclude "*.css" --exclude "*.png"

push-github:
	git branch -D gh-pages 2>/dev/null | true
	git branch -D gh-pages-draft 2>/dev/null | true
	git checkout -b gh-pages-draft && \
	git add -f build && \
	git commit -m "Deploy to gh-pages" && \
	git subtree split --prefix build -b gh-pages && \
	git push --force origin gh-pages:gh-pages && \
	git checkout -

cdn:
	curl -X DELETE "https://api.cloudflare.com/client/v4/zones/$(CLOUDFLARE_ZONE)/purge_cache" \
		-H "X-Auth-Email: $(CLOUDFLARE_EMAIL)" \
		-H "X-Auth-Key: $(CLOUDFLARE_CLIENT_API_KEY)" \
		-H "Content-Type: application/json" \
		--data '{"purge_everything":true}'

setup:
	brew update
	brew install wellington # wellington is used for sass building


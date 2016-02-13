run: clean static sass run-go
	
debug: clean static sass debug-go

build: clean static sass build-go

clean:
	rm -fR build

build-go:
	mkdir -p build
	go run *.go build

run-go:
	go run *.go server

debug-go:
	godebug run *.go server

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
	aws s3 sync build/ s3://readprayrepeat.com/ --acl public-read --content-type text/html --exclude "*.eot" --exclude "*.svg" --exclude "*.ttf" --exclude "*.woff" --exclude "*.js" --exclude "*.css" --exclude "*.png"
	aws s3 sync build/ s3://readprayrepeat.com/ --acl public-read --exclude "*" --include "*.eot" --include "*.svg" --include "*.ttf" --include "*.woff" --include "*.js" --include "*.css" --include "*.png"

push-github:
	git branch -D gh-pages 2>/dev/null | true
	git branch -D gh-pages-draft 2>/dev/null | true
	git checkout -b gh-pages-draft && \
	git add -f build && \
	git commit -m "Deploy to gh-pages" && \
	git subtree split --prefix build -b gh-pages && \
	git push --force origin gh-pages:gh-pages && \
	git checkout -

setup:
	brew update
	brew install wellington # wellington is used for sass building

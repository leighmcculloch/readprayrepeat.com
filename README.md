# Read, Pray, Repeat

Read and pray each day, using the reading plan my <a href="http://realitysf.com">church family</a> is using in 2016 for the <a href="http://bible.realitysf.com">Year of Biblical Literacy</a>.

## Usage

https://readprayrepeat.com

## About

### Why
I created this website to fulfill my own need to access the reading plan that my <a href="http://realitysf.com">church family</a> is using in 2016 for the <a href="http://bible.realitysf.com">Year of Biblical Literacy</a>. I needed the readings in plain text using an easy-to-read-aloud translation - void of ads, distractions, or verse numbers. I needed to open a single page at the beginning of my day and read, without interruption. This format is helping me stay focused and read the Bible.

### Reading Plan and Videos
The reading plan and videos were created by <a href="http://thebibleproject.tumblr.com/readscripture">The Bible Project</a>. This website is not associated with The Bible Project, and just displays the scripture and YouTube videos for each day in the reading plan.

## Development

### Setup

Requires `brew` to be installed.

```bash
go get github.com/leighmcculloch/readprayrepeat.com
make setup
```

### Using Locally

This will start a local server at `localhost:4567`.

```bash
make
```

### Deploying

This will build and deploy the website to an S3 bucket.

```bash
make build push
```

Alternatively, this will build and deploy the website to the repositories `origin/gh-pages`.

```bash
make build push-github
```

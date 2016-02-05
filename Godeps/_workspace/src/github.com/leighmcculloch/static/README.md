# static

A framework for generating static websites using Go and Go templates. It helps you use build static websites that are dynamically generated from sources like RSS feeds, databases, APIs, etc.

## Status

This framework was pulled out of a static website generated from Go. It's a work in progress towards being a standalone framework and is pretty not-configurable in it's current state.

## Install

```bash
go get github.com/leighmcculloch/static
```

## Usage

1. Create a `Static`

    ```go
    s := static.NewStatic()
    ```

2. Define a handler for each page that returns a data model, list of templates, and the entry point template to load for the page. The function will be called when building the page and the data model will be given to the template specific as the last return value. Templates are looked for in the source directory.

    ```go
    s.Handle("/index.html", func (path string) (data interface{}, tmpls []string, tmpl string) {
      ...
      return data, []string{"base.html", "index.html"}, "entry"
    }
    ```

3. (Optional) Configure the `Source` and `Build` directories, and other things.

    ```go
    s.SourceDir = "source" // The root directory for templates
    s.BuildDir = "build"   // The root directory for the built website
    s.ServerPort = 4567    // The port served on when running in server mode
    s.TemplateFuncs = template.FuncMap{...} // A map of functions available to templates
    ```

4. Call `Run`.

    ```go
    s.Run()
    ```

    A call to `Run` will build the website then exit if the first command line argument is `build`, otherwise it will serve the website on port `4567`.

    Instead of calling `Run` which uses the command line arguments, you can alternatively force one of the two actions to take place by calling `Build` or `Serve` directly.

## Simple Example

A simple example creating a page for each day of a year.

```go
package main

import "github.com/leighmcculloch/static"

const days = 365

func main() {
  s := static.NewStatic()

  s.Handle("/index.html", func(path string) (interface{}, []string, string) {
    return days, []string{"base.html", "index.html"}, "entry"
  })

  for d := 1; d <= days; d++ {
    day := struct { Number int } { d }

    path := fmt.Sprintf("/%d.html", day.Number)
    s.Handle(path, func(path string) (interface{}, []string, string) {
      return day, []string{"base.html", "day.html"}, "entry"
    })
  }

  s.Run()
}
```

## Typical Example

See [github.com/leighmcculloch/readprayrepeat.com](https://github.com/leighmcculloch/readprayrepeat.com).

## Others

Inspired by [middleman](https://middlemanapp.com/), which is a much easier to use and has many more features.

If your static website needs don't require writing code, checkout [Hugo](https://gohugo.io).

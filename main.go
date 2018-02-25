package main // import "github.com/leighmcculloch/today.bible"

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"4d63.com/biblepassageapi"
	"4d63.com/biblestats"
	"4d63.com/static"
)

const (
	cachePath = "cache"
)

var (
	biblesOrgAPIKey = func() string {
		apiKey := os.Getenv("API_KEY_BIBLESORG")
		if apiKey == "" {
			panic("API_KEY_BIBLESORG not set")
		}
		return apiKey
	}()
)

var (
	templateFuncs = template.FuncMap{
		"UnsafeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	templateIndex = template.Must(template.New("index").Funcs(templateFuncs).ParseFiles("source/base.html", "source/index.html"))
	templateDay   = template.Must(template.New("day").Funcs(templateFuncs).ParseFiles("source/base.html", "source/day.html"))
	templateVerse = template.Must(template.New("verse").Funcs(templateFuncs).ParseFiles("source/base.html", "source/verse.html"))
)

func main() {
	flagBuild := flag.Bool("build", false, "true to build, false to run server")
	flag.Parse()

	days := getDays()
	bible := biblepassageapi.Cache(
		biblepassageapi.NewBiblesOrg(biblesOrgAPIKey, biblepassageapi.GNT),
		cachePath,
	)

	mux := http.NewServeMux()
	paths := []string{}

	pages := []*pageDay{}
	var previousPage *pageDay
	for _, day := range days {
		page := &pageDay{
			PreviousPage: previousPage,
			Day:          day,
			Bible:        bible,
		}

		if previousPage != nil {
			previousPage.NextPage = page
		}
		previousPage = page

		pages = append(pages, page)
	}

	paths = append(paths, "/")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			templateIndex.ExecuteTemplate(w, "entry", pages)
		} else {
			fs := http.FileServer(http.Dir("build-assets"))
			fs.ServeHTTP(w, r)
		}
	})

	for i := range pages {
		page := pages[i]

		path := page.Path()
		paths = append(paths, path)
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			err := page.LoadPassages()
			if err != nil {
				log.Fatal(err)
			}
			err = templateDay.ExecuteTemplate(w, "entry", page)
			if err != nil {
				log.Fatal(err)
			}
		})
	}

	log.Printf("Registered %d handlers for reading", len(paths))

	var previousChapterPage *pagePassage
	// var previousVersePage *pagePassage
	for _, b := range biblestats.Books() {
		for c := 1; c <= biblestats.ChapterCount(b); c++ {
			page := &pagePassage{
				Reference:     fmt.Sprintf("%s %d", b, c),
				AbbrReference: fmt.Sprintf("%s %d", b.Abbr(), c),
				PreviousPage:  previousChapterPage,
				Bible:         bible,
			}
			if previousChapterPage != nil {
				previousChapterPage.NextPage = page
			}
			previousChapterPage = page

			path := page.Path()
			fmt.Println(path)
			paths = append(paths, path)
			mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				err := page.LoadPassages()
				if err != nil {
					log.Fatal(err)
				}
				err = templateVerse.ExecuteTemplate(w, "entry", page)
				if err != nil {
					log.Fatal(err)
				}
			})
		}
	}

	log.Printf("Registered %d handlers in total", len(paths))

	if *flagBuild {
		static.Build(static.DefaultOptions, mux, paths, func(e static.Event) {
			if e.Error == nil {
				return
			}
			log.Printf("%10s %d %-20s %s", e.Action, e.StatusCode, e.Path, e.Error)
		})
	} else {
		s := &http.Server{Addr: ":8080", Handler: mux}
		log.Fatal(s.ListenAndServe())
	}
}

func getDays() []Day {
	csvFile, err := os.Open("data/readings.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	return NewDaysFromCSV(csvFile)
}

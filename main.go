package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/leighmcculloch/static"
)

var (
	templateFuncs = template.FuncMap{
		"UnsafeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	templateIndex = template.Must(template.New("index").Funcs(templateFuncs).ParseFiles("source/base.html", "source/index.html"))
	templateDay   = template.Must(template.New("day").Funcs(templateFuncs).ParseFiles("source/base.html", "source/day.html"))
)

var (
	days   = getDays()
	bibles = []Bible{
		NewBiblesOrg(CEV),
		NewBiblesOrg(GNT),
		NewBibleNET(),
		NewESVAPI(),
	}
)

var build bool

func init() {
	flag.BoolVar(&build, "build", false, "true if build")
	flag.Parse()
}

func main() {
	mux := http.NewServeMux()
	paths := []string{}

	numberPages := len(bibles) * len(days)
	pages := make([]pageDay, 0, numberPages)
	uniqPages := pages[0:len(days)]
	for _, bible := range bibles {
		var previousPage *pageDay
		for _, day := range days {
			page := pageDay{
				PreviousPage: previousPage,
				Day:          day,
				Bible:        bible,
				Bibles:       bibles,
			}

			pages = append(pages, page)

			if previousPage != nil {
				previousPage.NextPage = &page
			}
			previousPage = &pages[len(pages)-1]
		}
	}

	paths = append(paths, "/")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			templateIndex.ExecuteTemplate(w, "entry", uniqPages)
		} else {
			fs := http.FileServer(http.Dir("build"))
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
			templateDay.ExecuteTemplate(w, "entry", page)
		})
	}

	log.Printf("Registered handlers for %d pages", len(pages))

	if build {
		static.Build(static.DefaultOptions(), mux, paths, func(e static.Event) {
			log.Println(e)
		})
	} else {
		s := &http.Server{Addr: ":8080", Handler: mux}
		log.Fatal(s.ListenAndServe())
	}
}

func getDays() []Day {
	csvFile, err := os.Open("data/data.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(records); i++ {
		record := records[i]
		for column := 0; column < len(record); column++ {
			record[column] = strings.Trim(record[column], " ")
			re := regexp.MustCompile(`\s+`)
			record[column] = re.ReplaceAllString(record[column], " ")
		}
	}

	return NewDaysFromCSV(records)
}

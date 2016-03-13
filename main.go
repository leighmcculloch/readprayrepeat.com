package main

import (
	"bufio"
	"encoding/csv"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/leighmcculloch/static"
	"github.com/leighmcculloch/static/build"
	"github.com/leighmcculloch/static/serve"
)

var (
	templateFuncs = template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToUpper":    strings.ToUpper,
		"UnsafeHTML": unsafeHTML,
	}
	templateIndex = template.Must(template.New("index").Funcs(templateFuncs).ParseFiles("source/base.html", "source/index.html"))
	templateDay   = template.Must(template.New("day").Funcs(templateFuncs).ParseFiles("source/base.html", "source/day.html"))
)

var (
	days   = loadDays()
	bibles = []Bible{
		NewBiblesOrg(CEV),
		NewBiblesOrg(GNT),
		NewBibleNET(),
		NewESVAPI(),
	}
)

func main() {
	s := static.New()

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

	s.AddPage("/index.html", func(w io.Writer, path string) error {
		return templateIndex.ExecuteTemplate(w, "entry", uniqPages)
	})

	for i := range pages {
		page := pages[i]

		path := page.Path()
		log.Printf("Registering handler for %s", path)
		s.AddPage(path, func(w io.Writer, path string) error {
			err := page.LoadPassages()
			if err != nil {
				log.Fatal(err)
			}
			return templateDay.ExecuteTemplate(w, "entry", page)
		})
	}

	var renderer static.Renderer
	switch os.Args[1] {
	case "build":
		builder := build.NewBuilder()
		builder.BuildConcurrency = 10
		renderer = builder
	default:
		renderer = serve.NewServer()
	}
	s.Render(renderer)
}

func loadDays() []Day {
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

	daysCount := len(records)
	days := make([]Day, daysCount)

	for i := 0; i < len(records); i++ {
		record := records[i]

		dayNumber := i + 1
		readingReference := record[0]
		prayerReference := record[1]
		watchYoutubeId := record[2]

		days[i] = Day{
			DayNumber:        dayNumber,
			ReadingReference: readingReference,
			PrayerReference:  prayerReference,
			WatchYoutubeId:   watchYoutubeId,
		}
	}

	return days
}

func unsafeHTML(s string) template.HTML {
	return template.HTML(s)
}

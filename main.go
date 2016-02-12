package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/leighmcculloch/static"
)

func main() {
	s := static.New()

	days := loadDays()

	s.Page("/index.html", func(path string) (interface{}, []string, string) {
		return days, []string{"base.html", "index.html"}, "entry"
	})

	bibles := []Bible{
		NewBiblesOrg(CEV),
		NewBiblesOrg(GNT),
		NewESVAPI(),
	}

	numberPages := len(bibles) * len(days)
	pages := make([]pageDay, 0, numberPages)
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

	for i := range pages {
		page := pages[i]

		path := fmt.Sprintf("%s.html", page.Path())
		log.Printf("Registering handler for %s", path)
		s.Page(path, func(path string) (interface{}, []string, string) {
			err := page.LoadPassages()
			if err != nil {
				log.Fatal(err)
			}
			return page, []string{"base.html", "day.html"}, "entry"
		})
	}

	s.Run()
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

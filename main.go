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
	s := static.NewStatic()

	days := loadDays()

	s.Handle("/index.html", func(path string) (interface{}, []string, string) {
		return days, []string{"base.html", "index.html"}, "entry"
	})

	bibles := []Bible{NewBiblesOrg(CEV)}

	for b, bible := range bibles {
		for _, day := range days {
			d := day

			var path string
			if b == 0 {
				path = fmt.Sprintf("/%d.html", d.DayNumber)
			} else {
				path = fmt.Sprintf("/%s/%d.html", bible.NameAbbr(), d.DayNumber)
			}
			s.Handle(path, func(path string) (interface{}, []string, string) {
				err := d.LoadPassages(bible)
				if err != nil {
					log.Fatal(err)
				}
				return d, []string{"base.html", "day.html"}, "entry"
			})
		}
	}

	s.Run()
}

func loadDays() []*Day {
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
	days := make([]*Day, daysCount)

	for i := 0; i < len(records); i++ {
		record := records[i]

		dayNumber := i + 1
		readingReference := record[0]
		prayerReference := record[1]
		watchYoutubeId := record[2]

		days[i] = NewDay(dayNumber, daysCount, readingReference, prayerReference, watchYoutubeId)
	}

	return days
}

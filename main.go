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

	s.Handle("/index.html", func(path string) (interface{}, []string) {
		return days, []string{"base.html", "index.html"}
	})

	for i := range days {
		d := days[i]
		p := fmt.Sprintf("/%d.html", d.DayNumber)
		s.Handle(p, func(path string) (interface{}, []string) {
			d.LoadPassages()
			return d, []string{"base.html", "day.html"}
		})
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

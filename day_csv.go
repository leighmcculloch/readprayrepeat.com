package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"4d63.com/collapsewhitespace"
)

func NewDaysFromCSV(r io.Reader) []Day {
	csvFile, err := os.Open("data/readings.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return NewDaysFromCSVRecords(records)
}

func NewDaysFromCSVRecords(records [][]string) []Day {
	daysCount := len(records)
	days := make([]Day, daysCount)

	for i := 0; i < len(records); i++ {
		record := records[i]

		dayNumber := i + 1

		days[i] = NewDayFromCSVRecord(dayNumber, record)
	}

	return days
}

func NewDayFromCSVRecord(dayNumber int, record []string) Day {
	for column := range record {
		record[column] = strings.Trim(record[column], " ")
		record[column] = collapsewhitespace.String(record[column])
	}
	return Day{
		DayNumber:      dayNumber,
		References:     []string{record[0], record[1]},
		WatchYoutubeId: record[2],
	}
}

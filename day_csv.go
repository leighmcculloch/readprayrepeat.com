package main

func NewDaysFromCSV(records [][]string) []Day {
	daysCount := len(records)
	days := make([]Day, daysCount)

	for i := 0; i < len(records); i++ {
		record := records[i]

		dayNumber := i + 1

		days[i] = NewDayFromCSV(dayNumber, record)
	}

	return days
}

func NewDayFromCSV(dayNumber int, record []string) Day {
	return Day{
		DayNumber:        dayNumber,
		ReadingReference: record[0],
		PrayerReference:  record[1],
		WatchYoutubeId:   record[2],
	}
}

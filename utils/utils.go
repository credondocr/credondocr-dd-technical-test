package utils

import (
	"strings"
	"time"
)

const DefaultLayout = "2006-01-02"

func GetDaysCountSince(from time.Time, until time.Time) int {
	days := 0
	day := from.Day()
	for from.Before(until) {
		from = from.Add(time.Hour * 24)
		nextDay := from.Day()
		if nextDay != day {
			days++
		}
		day = nextDay
	}
	return days
}

func GetDaysNotRequested(from time.Time, until time.Time, requested []string) []string {
	days := 0
	day := from.Day()
	notRequested := []string{}
	for from.Format(DefaultLayout) == until.Format(DefaultLayout) || from.Before(until) {
		if !Contains(requested, from.Format(DefaultLayout)) {
			notRequested = append(notRequested, from.Format(DefaultLayout))
		}
		from = from.Add(time.Hour * 24)
		nextDay := from.Day()
		if nextDay != day {
			days++
		}
		day = nextDay
	}
	return notRequested
}

func GetMonthsCountSince(from time.Time, until time.Time) int {
	months := 0
	month := from.Month()
	for from.Before(until) {
		from = from.Add(time.Hour * 24)
		nextMonth := from.Month()
		if nextMonth != month {
			months++
		}
		month = nextMonth
	}
	return months
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ParseStringToDate(date string) (time.Time, error) {
	timeParsed, err := time.Parse(DefaultLayout, strings.ReplaceAll(string(date), "\"", ""))
	if err != nil {
		return time.Now(), err
	}
	return timeParsed, nil
}

type CustomTime struct {
	time.Time
}

func (c *CustomTime) UnmarshalJSON(v []byte) error {
	var err error
	c.Time, err = time.Parse(DefaultLayout, strings.ReplaceAll(string(v), "\"", ""))
	if err != nil {
		return err
	}
	return nil
}

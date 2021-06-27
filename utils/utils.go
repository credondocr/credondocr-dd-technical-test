package utils

import (
	"time"

	"gorm.io/datatypes"
)

func GetDaysCountSince(from time.Time, to time.Time) int {
	days := 0
	day := from.Day()
	for from.Before(to) {
		from = from.Add(time.Hour * 24)
		nextDay := from.Day()
		if nextDay != day {
			days++
		}
		day = nextDay
	}
	return days
}

func GetDaysNotRequested(from time.Time, to time.Time, requested []string) []string {
	days := 0
	day := from.Day()
	notRequested := []string{}
	for from.Format("2006-01-02") == to.Format("2006-01-02") || from.Before(to) {
		if !Contains(requested, from.Format("2006-01-02")) {
			notRequested = append(notRequested, from.Format("2006-01-02"))
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

func GetMonthsCountSince(from time.Time, to time.Time) int {
	months := 0
	month := from.Month()
	for from.Before(to) {
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

type Data struct {
	ReleaseAt string         `json:"release_at"`
	Songs     datatypes.JSON `sql:"type:json" grom:"-" json:"songs"`
}

type Song struct {
	Artist string `json:"artist"`
	Name   string `json:"name"`
}

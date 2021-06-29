package services

import (
	"credondocr-dd-technical-test/dtos"
	"credondocr-dd-technical-test/models"
	"credondocr-dd-technical-test/utils"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/drgrib/iter"
	"github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const DailyLimitToBeRentable = 25

func ProcessDataToDataBase(db *gorm.DB, from time.Time, until time.Time) error {
	requestedDates, _ := models.GetRequestedDays(db, from, until)
	daysToProcess := utils.GetDaysNotRequested(from, until, requestedDates)

	if len(daysToProcess) < DailyLimitToBeRentable {
		requestDays(db, daysToProcess)
	} else {
		monthsToProcess := utils.GetMonthsCountSince(from, until)
		for range iter.N(monthsToProcess + 1) {
			d := 0
			if from.Day() > 1 {
				d = -from.Day() + 1
			}
			// d-1 means from the month to process, get the previous day from previous month
			// for example 2021-02-01  in line 44 will be 2021-01-30
			if from.AddDate(0, 1, d-1).After(until) {
				if utils.GetDaysCountSince(from, until) > DailyLimitToBeRentable {
					requestMonth(db, from)
				} else {
					requestedDates, err := models.GetRequestedDays(db, from, until)
					if err != nil {
						return err
					}
					pendingDaysToCached := utils.GetDaysNotRequested(from, until, requestedDates)
					requestDays(db, pendingDaysToCached)
				}
			} else {
				if utils.GetDaysCountSince(from, from.AddDate(0, 1, d-1)) > 25 {
					requestMonth(db, from)
				} else {
					requestedDates, err := models.GetRequestedDays(db, from, until)
					if err != nil {
						return err
					}
					pendingDaysToCached := utils.GetDaysNotRequested(from, from.AddDate(0, 1, d-1), requestedDates)
					requestDays(db, pendingDaysToCached)
				}
			}
			from = from.AddDate(0, 1, d)
		}
	}
	return nil
}

func ReadDataFromDatabase(db *gorm.DB, p dtos.Params) ([]dtos.Data, error) {
	return models.GetResponseWithFormat(db, p.From, p.Until, p.Artist)
}

func requestDays(db *gorm.DB, daysToProcess []string) {
	var waitgroup sync.WaitGroup
	waitgroup.Add(len(daysToProcess))
	for _, d := range daysToProcess {
		go func(day string) {
			res, err := getSingleDayRequest(day)
			log.Println("request day ", day)
			if err != nil {
				log.Fatal(err)
			}
			incomingSongs := []models.Song{}
			for _, row := range res {
				e, _ := json.Marshal(row.Stats)
				metadata := postgres.Jsonb{e}
				s := models.Song{ExternalId: row.SongId, Name: row.Name, ReleaseAt: datatypes.Date(row.ReleaseAt.Time), RequestedAt: datatypes.Date(time.Now()), Duration: row.Duration, Artist: row.Artist, MetaData: metadata}
				incomingSongs = append(incomingSongs, s)
			}
			_, err = models.StoreSongsIntoDatabase(db, incomingSongs)
			if err != nil {
				log.Fatal(err)
			}
			waitgroup.Done()
		}(d)
	}
	waitgroup.Wait()
}

func requestMonth(db *gorm.DB, from time.Time) {
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	go func() {
		res, err := getMonthRequest(from)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("request month ", from)
		incomingSongs := []models.Song{}
		for _, row := range res {
			e, _ := json.Marshal(row.Stats)
			metadata := postgres.Jsonb{e}
			s := models.Song{Name: row.Name, ReleaseAt: datatypes.Date(row.ReleaseAt.Time), RequestedAt: datatypes.Date(time.Now()), Duration: row.Duration, Artist: row.Artist, MetaData: metadata}
			incomingSongs = append(incomingSongs, s)
		}
		_, err = models.StoreSongsIntoDatabase(db, incomingSongs)
		if err != nil {
			log.Fatal(err)
		}
		waitgroup.Done()
	}()
	waitgroup.Wait()
}

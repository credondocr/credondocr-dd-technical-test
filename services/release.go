package services

import (
	"credondocr-dd-technical-test/models"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func FetchByDay(db *gorm.DB, daysToProcess []string) {
	var waitgroup sync.WaitGroup
	waitgroup.Add(len(daysToProcess))
	for _, d := range daysToProcess {
		go func(day string) {
			res, err := GetSingleDay(day)
			fmt.Println("request day ", day)
			if err != nil {
				log.Fatal(err)
			}
			incomingSongs := []models.Song{}
			for _, row := range res {
				e, _ := json.Marshal(row.Stats)
				metadata := postgres.Jsonb{e}
				s := models.Song{Name: row.Name, ReleaseAt: datatypes.Date(row.ReleaseAt.Time), RequestedAt: datatypes.Date(time.Now()), Duration: row.Duration, Artist: row.Artist, MetaData: metadata}
				incomingSongs = append(incomingSongs, s)
			}
			_, err = models.CreateSongs(db, incomingSongs)
			if err != nil {
				log.Fatal(err)
			}
			waitgroup.Done()
		}(d)
	}
	waitgroup.Wait()
}

func FetchByMonth(db *gorm.DB, from time.Time) {
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	go func() {
		res, err := GetMonth(from)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("request month ", from)
		incomingSongs := []models.Song{}
		for _, row := range res {
			e, _ := json.Marshal(row.Stats)
			metadata := postgres.Jsonb{e}
			s := models.Song{Name: row.Name, ReleaseAt: datatypes.Date(row.ReleaseAt.Time), RequestedAt: datatypes.Date(time.Now()), Duration: row.Duration, Artist: row.Artist, MetaData: metadata}
			incomingSongs = append(incomingSongs, s)
		}
		_, err = models.CreateSongs(db, incomingSongs)
		if err != nil {
			log.Fatal(err)
		}
		waitgroup.Done()
	}()
	waitgroup.Wait()
}

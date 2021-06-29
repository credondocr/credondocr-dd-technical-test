package models

import (
	"credondocr-dd-technical-test/dtos"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/myesui/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Song struct {
	Id          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ReleaseAt   datatypes.Date `json:"release_at"  time_format:"2006-01-02"`
	Duration    string         `json:"author"`
	Artist      string         `json:"artist"`
	Name        string         `json:"name"`
	RequestedAt datatypes.Date `gorm:"index" json:"requested_at" time_format:"2006-01-02" sql:"DEFAULT:current_timestamp"`
	MetaData    postgres.Jsonb `json:"meta_data"`
	ExternalId  string         `json:"external_id"`
}

type JsonTime struct {
	time.Time
}

func FindSongs(c *gin.Context, from string, to string, artist string) ([]Song, error) {
	db := c.MustGet("db").(*gorm.DB)
	sql := "release_at between ? and ?  AND requested_at <= now() + INTERVAL '30 DAYS' "
	var songs []Song
	if artist == "" {
		sql = sql + " AND artist = ?"
	}
	if err := db.Where(sql, from, to, artist).Order("release_at ASC").Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}

func GetResponseWithFormat(db *gorm.DB, from string, until string, artist string) ([]dtos.Data, error) {
	var results = []dtos.Data{}
	sql := `
		SELECT to_char(DATE (release_at)::date, 'YYYY-MM-DD') release_at, json_agg(json_build_object('name', name,'artist', artist)) Songs
			FROM songs
			where release_at between ? and ?  
			AND requested_at <= now() + INTERVAL '30 DAYS' `
	if artist != "" {
		sql = sql + `and artist = ? `
	}
	sql = sql + `group by release_at order by  release_at`
	err := db.Raw(sql, from, until, artist).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func StoreSongsIntoDatabase(db *gorm.DB, songs []Song) ([]Song, error) {
	if err := db.Create(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}

func GetRequestedDays(db *gorm.DB, from time.Time, to time.Time) ([]string, error) {
	var requestedDates []string
	sql := `select to_char(DATE (release_at)::date, 'YYYY-MM-DD')
					from songs
					where release_at between ? and ? and requested_at <= now() + INTERVAL '30 DAYS'
					group by release_at
	`
	if err := db.Raw(sql, from, to).Scan(&requestedDates).Error; err != nil {
		return nil, err
	}
	return requestedDates, nil
}

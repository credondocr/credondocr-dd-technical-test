package models

import (
	"credondocr-dd-technical-test/utils"
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
}

type JsonTime struct {
	time.Time
}

func FindSongs(c *gin.Context, from time.Time, to time.Time, artist string) ([]Song, error) {
	db := c.MustGet("db").(*gorm.DB)
	var songs []Song
	if artist == "" {
		if err := db.Where("release_at between ? and ?  AND requested_at <= now() + INTERVAL '30 DAYS'", from, to).Order("release_at ASC").Find(&songs).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Where("release_at between ? and ?  AND requested_at <= now() + INTERVAL '30 DAYS' AND artist = ?", from, to, artist).Order("release_at ASC").Find(&songs).Error; err != nil {
			return nil, err
		}
	}

	return songs, nil
}

func GetResult(c *gin.Context, from time.Time, to time.Time, artist string) ([]utils.Data, error) {
	db := c.MustGet("db").(*gorm.DB)
	var results = []utils.Data{}
	if artist == "" {
		sql := `
		SELECT to_char(DATE (release_at)::date, 'YYYY-MM-DD') release_at, json_agg(json_build_object('name', name,'artist', artist)) Songs
			FROM songs
			where release_at between ? and ?  AND requested_at <= now() + INTERVAL '30 DAYS' 
			group by release_at 
			order by  release_at`

		_ = db.Raw(sql, from, to).Scan(&results).Error
	} else {
		sql := `
		SELECT to_char(DATE (release_at)::date, 'YYYY-MM-DD') release_at, json_agg(json_build_object('name', name,'artist', artist)) Songs
			FROM songs
			where release_at between ? and ?  AND requested_at <= now() + INTERVAL '30 DAYS' and artist = ?
			group by release_at 
			order by  release_at`
		_ = db.Raw(sql, from, to, artist).Scan(&results).Error
	}
	return results, nil
}

func CreateSongs(c *gin.Context, songs []Song) ([]Song, error) {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}

func GetValidDays(c *gin.Context, from time.Time, to time.Time) ([]string, error) {
	db := c.MustGet("db").(*gorm.DB)
	var validDates []string
	if err := db.Raw("select to_char(DATE (release_at)::date, 'YYYY-MM-DD') from songs where release_at between ? and ? and requested_at <= now() + INTERVAL '30 DAYS'  group by release_at", from, to).Scan(&validDates).Error; err != nil {
		return nil, err
	}
	return validDates, nil
}

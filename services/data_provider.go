package services

import (
	"credondocr-dd-technical-test/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Data struct {
	SongId    string     `json:"song_id"`
	ReleaseAt CustomTime `json:"released_at" time_format:"2006-01-02"`
	Duration  string     `json:"duration"`
	Artist    string     `json:"artist"`
	Name      string     `json:"name"`
	Stats     Stats      `json:"stats"`
}

type Stats struct {
	LastPlayedAt  int `json:"last_played_at"`
	TimesPlayedAt int `json:"times_played_at"`
	GlobalRank    int `json:"global_rank"`
}

type CustomTime struct {
	time.Time
}

func (c *CustomTime) UnmarshalJSON(v []byte) error {
	var err error
	c.Time, err = time.Parse("2006-01-02", strings.ReplaceAll(string(v), "\"", ""))
	if err != nil {
		return err
	}
	return nil
}

func GetSingleDay(date string) ([]Data, error) {
	config := config.GetConfig()
	resp, err := http.Get(fmt.Sprintf("%s/daily?api_key=ec093dd5-bbe3-4d8e-bdac-314b40afb796&released_at=%s", config.GetString("provider.uri"), date))
	if err != nil {
		return nil, err
	}

	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		log.Fatal(readErr)
		return nil, err
	}
	requested := []Data{}
	jsonErr := json.Unmarshal(body, &requested)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return requested, nil
}

func GetMonth(from time.Time) ([]Data, error) {
	config := config.GetConfig()
	resp, err := http.Get(fmt.Sprintf("%s/monthly?api_key=ec093dd5-bbe3-4d8e-bdac-314b40afb796&released_at=%s", config.GetString("provider.uri"), from.Format("2006-01")))
	if err != nil {
		return nil, err
	}

	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		log.Fatal(readErr)
		return nil, err
	}
	requested := []Data{}
	jsonErr := json.Unmarshal(body, &requested)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return requested, nil
}

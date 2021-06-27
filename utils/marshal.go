package utils

import (
	"strings"
	"time"
)

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

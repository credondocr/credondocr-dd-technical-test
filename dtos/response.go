package dtos

import "gorm.io/datatypes"

type Data struct {
	ReleaseAt string         `json:"release_at"`
	Songs     datatypes.JSON `sql:"type:json" grom:"-" json:"songs"`
}

type Song struct {
	Artist string `json:"artist"`
	Name   string `json:"name"`
}
package controllers

import (
	"credondocr-dd-technical-test/dtos"
	"credondocr-dd-technical-test/services"
	"credondocr-dd-technical-test/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type ReleaseController struct{}

func (h ReleaseController) Releases(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var p dtos.Params

	if err := c.ShouldBindWith(&p, binding.Query); err != nil {
		triggerDefaultError(c, err, err.Error())
		return
	}

	from, err := utils.ParseStringToDate(p.From)

	if err != nil {
		triggerDefaultError(c, err, "Invalid from dates, please check the params")
		return
	}

	until, err := utils.ParseStringToDate(p.Until)

	if err != nil {
		triggerDefaultError(c, err, "Invalid until dates, please check the params")
		return
	}
	if from.After(until) {
		triggerDefaultError(c, err, "Invalid dates, please check the params")
		return
	}

	if err := services.ProcessDataToDataBase(db, from, until); err != nil {
		triggerDefaultError(c, err, "Error reading the database")
	}

	result, err := services.ReadDataFromDatabase(db, p)

	if err != nil {
		triggerDefaultError(c, err, "Error reading the database")
		return
	}

	c.JSON(http.StatusOK, result)
}

func triggerDefaultError(c *gin.Context, err error, message string) {
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": message})
	}
}

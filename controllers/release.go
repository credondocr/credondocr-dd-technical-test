package controllers

import (
	"credondocr-dd-technical-test/models"
	"credondocr-dd-technical-test/services"
	"credondocr-dd-technical-test/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/drgrib/iter"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type ReleaseController struct{}

type Params struct {
	From   string `form:"from" binding:"required"`
	Until  string `form:"until" binding:"required"`
	Artist string `form:"artist"`
}

func (h ReleaseController) Releases(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var p Params
	if err := c.ShouldBindWith(&p, binding.Query); err == nil {
		from, err := time.Parse("2006-01-02", strings.ReplaceAll(string(p.From), "\"", ""))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "Invalid from dates, please check the params"})
			return
		}
		until, err := time.Parse("2006-01-02", strings.ReplaceAll(string(p.Until), "\"", ""))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "Invalid until dates, please check the params"})
			return
		}
		if from.After(until) {
			c.JSON(http.StatusOK, gin.H{"error": "Invalid dates, please check the params"})
			return
		}
		requestedDates, _ := models.GetRequestedDays(db, from, until)
		daysToProcess := utils.GetDaysNotRequested(from, until, requestedDates)

		if len(daysToProcess) < 25 {
			services.FetchByDay(db, daysToProcess)
		} else {
			f := utils.GetMonthsCountSince(from, until)
			for range iter.N(f + 1) {
				d := 0
				if from.Day() > 1 {
					d = -from.Day() + 1
				}
				if from.AddDate(0, 1, d-1).After(until) {
					if utils.GetDaysCountSince(from, until) > 25 {
						services.FetchByMonth(db, from)
					} else {
						requestedDates, _ := models.GetRequestedDays(db, from, until)
						services.FetchByDay(db, utils.GetDaysNotRequested(from, until, requestedDates))
					}
				} else {
					if utils.GetDaysCountSince(from, from.AddDate(0, 1, d-1)) > 25 {
						services.FetchByMonth(db, from)
					} else {
						requestedDates, _ := models.GetRequestedDays(db, from, until)
						services.FetchByDay(db, utils.GetDaysNotRequested(from, from.AddDate(0, 1, d-1), requestedDates))
					}
				}
				from = from.AddDate(0, 1, d)
			}
		}
		// time.Sleep(time.Duration(15 * time.Second))
		result, _ := models.GetResult(c.MustGet("db").(*gorm.DB), from, until, p.Artist)
		fmt.Println(result)
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

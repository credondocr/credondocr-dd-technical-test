package controllers

import (
	"credondocr-dd-technical-test/models"
	"credondocr-dd-technical-test/services"
	"credondocr-dd-technical-test/utils"
	"net/http"
	"strings"
	"time"

	"github.com/drgrib/iter"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ReleaseController struct{}

type Params struct {
	From   string `form:"from" binding:"required"`
	Until  string `form:"until" binding:"required"`
	Artist string `form:"artist"`
}

func (h ReleaseController) Releases(c *gin.Context) {

	var b Params
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		// if time.Parse("2006-01-02", b.From).After("2006-01-02", time.Parse(b.Until)) {
		// 	c.JSON(http.StatusOK, gin.H{"error": "Invalid dates, please check the params"})
		// 	return
		// }
		from, _ := time.Parse("2006-01-02", strings.ReplaceAll(string(b.From), "\"", ""))
		until, _ := time.Parse("2006-01-02", strings.ReplaceAll(string(b.Until), "\"", ""))
		validDates, _ := models.GetValidDays(c, from, until)

		daysToProcess := utils.GetDaysNotRequested(from, until, validDates)

		if len(daysToProcess) < 25 {
			services.FetchByDay(c, daysToProcess)
		} else {
			f := utils.GetMonthsCountSince(from, until)
			for range iter.N(f + 1) {
				d := 0
				if from.Day() > 1 {
					d = -from.Day() + 1
				}
				if from.AddDate(0, 1, d-1).After(until) {
					if utils.GetDaysCountSince(from, until) > 25 {
						services.FetchByMonth(c, from)
					} else {
						validDates, _ := models.GetValidDays(c, from, until)
						services.FetchByDay(c, utils.GetDaysNotRequested(from, until, validDates))
					}
				} else {
					if utils.GetDaysCountSince(from, from.AddDate(0, 1, d-1)) > 25 {
						services.FetchByMonth(c, from)
					} else {
						validDates, _ := models.GetValidDays(c, from, until)
						services.FetchByDay(c, utils.GetDaysNotRequested(from, from.AddDate(0, 1, d-1), validDates))
					}
				}
				from = from.AddDate(0, 1, d)
			}
		}

		result, _ := models.GetResult(c, from, until, b.Artist)
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

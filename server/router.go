package server

import (
	"credondocr-dd-technical-test/controllers"
	"credondocr-dd-technical-test/models"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	db := models.SetupModels()
	// Provide db variable to controllers
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	health := new(controllers.HealthController)
	release := new(controllers.ReleaseController)

	r.GET("/health", health.Status)
	r.GET("/releases", release.Releases)

	return r

}

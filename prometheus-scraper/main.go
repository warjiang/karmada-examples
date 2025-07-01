package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"prometheus-scraper/pkg/api"
	"prometheus-scraper/pkg/database"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	api.RegisterRoutes(r, db)

	r.GET("/", func(c *gin.Context) {
		c.File("web/index.html")
	})
	r.Static("/static", "./web/static")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

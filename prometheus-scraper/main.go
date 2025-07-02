package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"prometheus-scraper/pkg/api"
	"prometheus-scraper/pkg/config"
	"prometheus-scraper/pkg/database"
	"prometheus-scraper/pkg/scraper"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Printf("failed to load config, using default: %v", err)
		cfg = &config.Config{
			IntervalSeconds: 60,
			Targets: []config.ScrapeTarget{
				{
					Namespace: "karmada-system",
					Name:      "karmada-scheduler-7bd4659f9f-8lfb5",
					Port:      "10351",
				},
			},
		}
	}

	go startScraping(cfg, db)

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

func loadConfig(filename string) (*config.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func startScraping(cfg *config.Config, db *database.DB) {
	for {
		log.Println("Starting periodic scrape...")
		for _, target := range cfg.Targets {
			metrics, err := scraper.ScrapeFromPod(target.Namespace, target.Name, target.Port)
			if err != nil {
				log.Printf("Error scraping %s/%s:%s: %v", target.Namespace, target.Name, target.Port, err)
				continue
			}
			for _, metric := range metrics {
				metric.Timestamp = time.Now().Unix()
				if err := db.InsertMetric(metric); err != nil {
					log.Printf("Error inserting metric %s: %v", metric.Name, err)
				}
			}
			log.Printf("Scraped %d metrics from %s/%s:%s", len(metrics), target.Namespace, target.Name, target.Port)
		}
		time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
	}
}

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"prometheus-scraper/pkg/database"
	"prometheus-scraper/pkg/scraper"
)

func RegisterRoutes(r *gin.Engine, db *database.DB) {
	r.GET("/metrics/:name", getMetrics(db))
	r.GET("/metrics/names", getMetricNames(db))
	r.POST("/scrape", scrapeMetrics(db))
	r.POST("/scrape/pod", scrapePodMetrics(db))
}

func getMetrics(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		labels := make(map[string]string)
		for k, v := range c.Request.URL.Query() {
			if k != "name" {
				labels[k] = v[0]
			}
		}
		metrics, err := db.GetMetrics(name, labels)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, metrics)
	}
}

func getMetricNames(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		names, err := db.GetMetricNames()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, names)
	}
}

func scrapeMetrics(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type scrapeRequest struct {
			URL string `json:"url"`
		}
		req := new(scrapeRequest)
		if err := c.BindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		metrics, err := scraper.Scrape(req.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, metric := range metrics {
			if err := db.InsertMetric(metric); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, nil)
	}
}

func scrapePodMetrics(db *database.DB) gin.HandlerFunc {
	/*
		curl -X POST 'http://localhost:8080/scrape/pod' -d '{"namespace":"karmada-system","name": "karmada-scheduler-7bd4659f9f-8lfb5", "port": "10351"}'
	*/
	return func(c *gin.Context) {
		type scrapePodRequest struct {
			Namespace string `json:"namespace"`
			Name      string `json:"name"`
			Port      string `json:"port"`
		}
		req := new(scrapePodRequest)
		if err := c.BindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		metrics, err := scraper.ScrapeFromPod(req.Namespace, req.Name, req.Port)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, metric := range metrics {
			if err := db.InsertMetric(metric); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, nil)
	}
}

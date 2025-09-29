package handlers

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"pikalink-backend/database"
	"pikalink-backend/utils"
	"github.com/gin-gonic/gin"
)

func GetAnalysisMonths(c *gin.Context) {
	logDir := utils.GetLogsDir()
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read log directory"})
		return
	}

	var months []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".db" {
			month := strings.TrimSuffix(file.Name(), ".db")
			months = append(months, month)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(months)))

	c.JSON(http.StatusOK, months)
}

func AnalyzeLogs(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Month query parameter is required"})
		return
	}

	shortURL := c.Query("short_url")

	analysis, err := database.GetAnalysis(month, shortURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze logs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(c *gin.Context) {
	log.Println("Request received:", c.Request.Method, c.Request.URL.Path)
	c.Next()
}

func ValidatePostRequest(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read request body"})
		c.Abort()
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var request struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON formate"})
		c.Abort()
		return
	}

	if request.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		c.Abort()
		return
	}

	parsedURL, err := url.ParseRequestURI(request.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		c.Abort()
		return
	}
	c.Next()
}

func ValidateGetRequest(c *gin.Context) {
	shortURL := c.Param("shortUrl")

	var validShortURL = regexp.MustCompile(`^[a-zA-Z0-9_]{10}$`)

	if !validShortURL.MatchString(shortURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid short URL format"})
		c.Abort()
		return
	}
	c.Next()
}

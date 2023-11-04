package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	req := gin.Default()
	m := make(map[string]string)

	req.GET("/ping", func(context *gin.Context) {
		url := context.Query("url")

		if url == "" {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'url' parameter"})
			return
		}

		shortIdentifier := generateShortIdentifier(url)

		m[shortIdentifier] = url

		context.JSON(http.StatusOK, gin.H{"short_url": "localhost:8080/" + shortIdentifier})
	})

	req.GET("/pong", func(context *gin.Context) {
		shortUrl := context.Query("short_url")
		parts := strings.Split(shortUrl, "/")
		hash := parts[len(parts)-1]
		value, exists := m[hash]

		if shortUrl == "" || !exists {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'short_url' parameter"})
			return
		}

		context.JSON(http.StatusOK, gin.H{"url": value})

	})

	req.Run(":8080")
}

func generateShortIdentifier(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))[:8] // Take the first 8 characters of the hash
}

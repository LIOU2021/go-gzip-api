package main

import (
	gzipOrigin "compress/gzip"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression), GzipDecompress)
	r.POST("/data", func(c *gin.Context) {
		b, _ := io.ReadAll(c.Request.Body)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, world! from gin",
			"request": string(b),
		})
	})

	if err := r.Run("127.0.0.1:8080"); err != nil {
		log.Fatal(err)
	}
}

func GzipDecompress(c *gin.Context) {
	if c.GetHeader("Content-Encoding") == "gzip" {
		reader, err := gzipOrigin.NewReader(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip format"})
			return
		}
		defer reader.Close()
		c.Request.Body = http.MaxBytesReader(c.Writer, reader, c.Request.ContentLength)
	}
}

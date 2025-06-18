package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// SimpleRequestLogger returns a basic gin middleware that logs requests
func SimpleRequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format with emojis and detailed information
		statusEmoji := "✅"
		if param.StatusCode >= 400 {
			statusEmoji = "❌"
		} else if param.StatusCode >= 300 {
			statusEmoji = "🔄"
		}

		return fmt.Sprintf("%s [%s] %s %s %d %s \"%s\" %s \"%s\" %s\n",
			statusEmoji,
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			param.StatusCode,
			param.Latency,
			param.Path,
			param.Request.Proto,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// DetailedRequestLogger logs request and response details for debugging
func DetailedRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request details
		log.Printf("📥 REQUEST: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		log.Printf("📋 Headers: %+v", c.Request.Header)

		// Log request body for POST/PUT requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.Body != nil {
				bodyBytes, _ := io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Try to parse as JSON for pretty printing
				var jsonBody interface{}
				if json.Unmarshal(bodyBytes, &jsonBody) == nil {
					prettyBody, _ := json.MarshalIndent(jsonBody, "", "  ")
					log.Printf("📄 Request Body:\n%s", string(prettyBody))
				} else {
					log.Printf("📄 Request Body: %s", string(bodyBytes))
				}
			}
		}

		// Capture start time
		start := time.Now()

		// Process request
		c.Next()

		// Log response details
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		statusEmoji := "✅"
		if statusCode >= 400 {
			statusEmoji = "❌"
		} else if statusCode >= 300 {
			statusEmoji = "🔄"
		}

		log.Printf("%s RESPONSE: %d in %v", statusEmoji, statusCode, latency)

		// Log errors if any
		if len(c.Errors) > 0 {
			log.Printf("🚨 Errors: %+v", c.Errors)
		}

		log.Printf("─────────────────────────────────────")
	}
}

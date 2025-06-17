package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs all HTTP requests and responses
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RequestResponseLogger logs detailed request and response information
func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request
		log.Printf("=== REQUEST ===")
		log.Printf("Method: %s", c.Request.Method)
		log.Printf("URL: %s", c.Request.URL.String())
		log.Printf("Headers: %v", c.Request.Header)

		// Read and log request body if present
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bodyBytes) > 0 {
				log.Printf("Request Body: %s", string(bodyBytes))
				// Restore the body for further processing
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Create a response writer wrapper to capture response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// Log response
		log.Printf("=== RESPONSE ===")
		log.Printf("Status: %d", c.Writer.Status())
		log.Printf("Latency: %v", latency)
		log.Printf("Response Headers: %v", c.Writer.Header())

		// Log response body
		responseBody := blw.body.String()
		if responseBody != "" {
			log.Printf("Response Body: %s", responseBody)
		}

		// Log any errors
		if len(c.Errors) > 0 {
			log.Printf("Errors: %v", c.Errors)
		}

		log.Printf("================")
	}
}

// bodyLogWriter wraps gin.ResponseWriter to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ErrorHandler middleware for centralized error handling
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Request error: %v", err)

			// If no response has been written yet, send error response
			if !c.Writer.Written() {
				c.JSON(500, gin.H{
					"error":   "Internal server error",
					"details": err.Error(),
				})
			}
		}
	}
}

package logger

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Собираем ошибки из контекста, если они есть
		var errorMessages []string
		for _, e := range c.Errors {
			errorMessages = append(errorMessages, e.Error())
		}

		if len(errorMessages) > 0 {
			log.Printf("[ERROR] %d %s %s in %v | Errors: %s", statusCode, method, path, duration, strings.Join(errorMessages, "; "))
		} else {
			log.Printf("[INFO] %d %s %s in %v", statusCode, method, path, duration)
		}
	}
}

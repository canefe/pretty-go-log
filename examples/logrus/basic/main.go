package main

import (
	"github.com/canefe/pretty-go-log/logrus/pretty"
)

func main() {
	// Create a new logger with default settings
	log := pretty.New()

	// Log messages at different levels
	log.Debug("This is a debug message")
	log.Info("Application started successfully")
	log.Info("[Server] Listening on port 8080")
	log.Warn("This is a warning message")
	log.Error("This is an error message")
	log.Trace("hey")

	// With fields - fields appear in dim color at the end of the line
	log.WithField("user", "alice").Info("User logged in")
	log.WithField("duration", "150ms").Info("[Request] GET /api/users")

	// Multiple fields - automatically sorted alphabetically
	log.WithFields(map[string]interface{}{
		"user_id":  123,
		"endpoint": "/api/data",
		"method":   "POST",
		"status":   200,
	}).Info("[API] Request completed")
}

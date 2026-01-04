package main

import (
	"github.com/canefe/pretty-go-log/logrus/pretty"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a logger with JSON output format
	log := pretty.New(
		pretty.WithLevel(logrus.InfoLevel),
		pretty.WithOutput(pretty.OutputConsole),
		pretty.WithFormat(pretty.FormatJSON),
	)

	log.Info("Application started")
	log.WithField("version", "1.0.0").Info("Build information")
	log.WithFields(map[string]interface{}{
		"user_id": 123,
		"action":  "login",
		"ip":      "192.168.1.1",
	}).Info("User activity")

	log.WithField("error", "connection timeout").Error("Database connection failed")
}

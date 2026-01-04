package main

import (
	"github.com/canefe/pretty-go-log/logrus/pretty"
	"github.com/sirupsen/logrus"
)

func main() {
	log := pretty.New(
		pretty.WithLevel(logrus.DebugLevel),
		pretty.WithOutput(pretty.OutputMulti), // Output to both console and file
		pretty.WithFormat(pretty.FormatPlain),
		pretty.WithNamespace("App"),
		pretty.WithFile("logs/service.log"),
	)

	log.Debug("[Init] Starting application with multi-output logging")
	log.Info("[Server] Server initialized")
	log.Info("[Database] Connected to PostgreSQL")
	log.Warn("[Cache] Redis connection slow")
	log.Error("[API] External API timeout")

	log.WithFields(logrus.Fields{
		"user_id":  456,
		"endpoint": "/api/data",
		"method":   "POST",
		"duration": "250ms",
	}).Info("[Request] API request completed")

	log.Info("Check both console output and logs/service.log")
}

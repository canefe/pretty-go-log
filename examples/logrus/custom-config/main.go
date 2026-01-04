package main

import (
	"github.com/canefe/pretty-go-log/logrus/pretty"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a logger with custom configuration
	log := pretty.New(
		pretty.WithLevel(logrus.DebugLevel),
		pretty.WithOutput(pretty.OutputConsole),
		pretty.WithFormat(pretty.FormatPlain),
	)

	log.Debug("[Init] Initializing application")
	log.Info("[Config] Configuration loaded")
	log.Info("[Database] Connected to database")
	log.Warn("[Cache] Cache miss for key: user:123")
	log.Error("[API] Failed to connect to external service")

	// Some other logger
	ordersLog := pretty.New(
		pretty.WithLevel(logrus.DebugLevel),
		pretty.WithOutput(pretty.OutputMulti), // multi
		pretty.WithFormat(pretty.FormatPlain),
		pretty.WithNamespace("Orders"),
		pretty.WithFile("logs/orders.log"),
	)

	ordersLog.WithFields(logrus.Fields{
		"id":     "123",
		"userId": "12345",
	}).Info("[Order] New order!")
}

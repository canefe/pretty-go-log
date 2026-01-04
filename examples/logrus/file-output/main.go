package main

import (
	"github.com/canefe/pretty-go-log/logrus/pretty"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a logger that writes to a file
	// The file will be rotated automatically using lumberjack
	log := pretty.New(
		pretty.WithLevel(logrus.InfoLevel),
		pretty.WithOutput(pretty.OutputFile), // Output to file (logs/service.log by default)
		pretty.WithFormat(pretty.FormatPlain),
		pretty.WithFile("logs/service.log"),
		pretty.WithoutCaller(),
	)

	log.Info("Application started - this will be written to file")
	log.Info("[Server] Server is running")

	for i := 0; i < 10; i++ {
		log.WithField("iteration", i).Info("[Task] Processing item")
	}

	log.Info("Application finished - check logs/service.log")
}

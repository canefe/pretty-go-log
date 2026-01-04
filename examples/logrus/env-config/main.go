package main

import (
	"os"

	"github.com/canefe/pretty-go-log/logrus/pretty"
)

func main() {
	// The logger can be configured via environment variables
	// Set these before running:
	// export LOG_LEVEL=debug
	// export LOG_OUTPUT=console
	// export LOG_FORMAT=plain

	// When options are omitted, the logger will
	// check the corresponding environment variables
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_OUTPUT", "console")
	os.Setenv("LOG_FORMAT", "plain")

	log := pretty.New()

	log.Debug("[Env] Logger configured via environment variables")
	log.Info("[Env] LOG_LEVEL=" + os.Getenv("LOG_LEVEL"))
	log.Info("[Env] LOG_OUTPUT=" + os.Getenv("LOG_OUTPUT"))
	log.Info("[Env] LOG_FORMAT=" + os.Getenv("LOG_FORMAT"))
}

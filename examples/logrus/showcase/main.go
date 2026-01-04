package main

import (
	"github.com/canefe/pretty-go-log/logrus/pretty"
	"github.com/sirupsen/logrus"
)

func main() {
	// Demo 1: StyleDefault - Standard left-aligned brackets
	println("=== StyleDefault (Standard Left-Aligned) ===\n")
	defaultLogger := logrus.New()
	defaultLogger.SetLevel(logrus.DebugLevel)
	defaultLogger.SetReportCaller(true)
	defaultLogger.SetFormatter(pretty.NewCustomFormatter(
		pretty.WithColors(true),
		pretty.WithCaller(true, logrus.WarnLevel),
		pretty.WithBracketPadding(16),
		pretty.WithColorBrackets(true),
		pretty.WithTagStyle(pretty.StyleDefault, "·"),
	))

	defaultLogger.Debug("[Init] Initializing application")
	defaultLogger.Info("[Server] Starting HTTP server on port 8080")
	defaultLogger.Info("[Database] Connection pool established")
	defaultLogger.WithFields(logrus.Fields{
		"host": "localhost",
		"port": 5432,
	}).Info("[DB] Connected to PostgreSQL")
	defaultLogger.Warn("[API] Rate limit approaching threshold")
	defaultLogger.WithFields(logrus.Fields{
		"user":  "admin",
		"error": "invalid token",
	}).Error("[Auth] Authentication failed")

	println("\n=== StyleCenter (Centered with Decorative Padding) ===\n")
	// Demo 2: StyleCenter - Centered brackets with decorative padding
	centeredLogger := logrus.New()
	centeredLogger.SetLevel(logrus.DebugLevel)
	centeredLogger.SetReportCaller(true)
	centeredLogger.SetFormatter(pretty.NewCustomFormatter(
		pretty.WithColors(true),
		pretty.WithCaller(true, logrus.WarnLevel),
		pretty.WithBracketPadding(16),
		pretty.WithColorBrackets(true),
		pretty.WithTagStyle(pretty.StyleCenter, " "),
	))

	centeredLogger.Debug("[Init] Initializing application")
	centeredLogger.Info("[Server] Starting HTTP server on port 8080")
	centeredLogger.Info("[Database] Connection pool established")
	centeredLogger.WithFields(logrus.Fields{
		"host": "localhost",
		"port": 5432,
	}).Info("[DB] Connected to PostgreSQL")
	centeredLogger.Warn("[API] Rate limit approaching threshold")
	centeredLogger.WithFields(logrus.Fields{
		"user":  "admin",
		"error": "invalid token",
	}).Error("[Auth] Authentication failed")

	println("\n=== Comparison: Different Padding Characters ===\n")
	// Demo 4: Different padding characters with StyleCenter
	chars := []struct {
		char string
		name string
	}{
		{"=", "Equals"},
		{"-", "Dash"},
		{"·", "Middle Dot"},
		{"•", "Bullet"},
	}

	for _, ch := range chars {
		println("--- " + ch.name + " (" + ch.char + ") ---")
		logger := logrus.New()
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(pretty.NewCustomFormatter(
			pretty.WithColors(true),
			pretty.WithBracketPadding(16),
			pretty.WithColorBrackets(true),
			pretty.WithTagStyle(pretty.StyleCenter, ch.char),
		))
		logger.Info("[Init] Application started")
		logger.Info("[Server] Listening on :8080")
		println()
	}

	println("=== Visual Alignment Demo ===\n")
	println("Notice how all messages start at the SAME column position!")
	println("This makes scanning logs 2x faster.\n")

	alignLogger := logrus.New()
	alignLogger.SetLevel(logrus.InfoLevel)
	alignLogger.SetReportCaller(true)
	alignLogger.SetFormatter(pretty.NewCustomFormatter(
		pretty.WithColors(true),
		pretty.WithCaller(true, logrus.WarnLevel),
		pretty.WithBracketPadding(16),
		pretty.WithColorBrackets(true),
		pretty.WithTagStyle(pretty.StyleCenter, " "),
	))

	// Different tag lengths - all perfectly aligned!
	alignLogger.Info("[DB] Short tag")
	alignLogger.Info("[Server] Medium tag")
	alignLogger.Info("[Application] Longer tag")
	alignLogger.Info("[Initialize] Even longer tag")
	alignLogger.WithFields(logrus.Fields{
		"duration": "125ms",
		"status":   "success",
	}).Info("[API] With structured fields")
	alignLogger.Warn("[Cache] Warning with caller info shown below")
	alignLogger.Error("[Critical] Error with caller info shown below")

	println("\n=== StyleRight (Left-Aligned with Trailing Padding) ===\n")

	// Demo 3: StyleRight - Left-aligned brackets with trailing padding
	formatter := pretty.NewCustomFormatter(
		pretty.WithColors(true),
		pretty.WithCaller(true, logrus.WarnLevel),
		pretty.WithBracketPadding(10),
		pretty.WithColorBrackets(true),
		pretty.WithTagStyle(pretty.StyleRight, "_"),
	)

	rightLogger := pretty.New(
		pretty.WithLevel(logrus.DebugLevel),
		pretty.WithOutput(pretty.OutputMulti),
		pretty.WithFormat(pretty.FormatPlain),
		pretty.WithCustomFormat(*formatter),
		pretty.WithFile("logs/service.log"),
	)

	rightLogger.Debug("[Init] Initializing application")
	rightLogger.Info("[Server] Starting HTTP server on port 8080")
	rightLogger.Info("[Database] Connection pool established")
	rightLogger.WithFields(logrus.Fields{
		"host": "localhost",
		"port": 5432,
	}).Info("[DB] Connected to PostgreSQL")
	rightLogger.Warn("[API] Rate limit approaching threshold")
	rightLogger.WithFields(logrus.Fields{
		"user":  "admin",
		"error": "invalid token",
	}).Error("[Auth] Authentication failed")

	println("\n=== Three Styles Side-by-Side ===\n")
	println("StyleDefault:  [Auth]         Clean and minimal")
	println("StyleCenter:   ...[Auth]...   Balanced and decorative")
	println("StyleRight:    [Auth].......  Modern and distinct\n")
}

package pretty

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type FormatType int

const (
	FormatRaw   FormatType = iota // 0
	FormatPlain                   // 1
	FormatJSON                    // 2
)

type OutputType int

const (
	OutputConsole OutputType = iota
	OutputFile
	OutputMulti
)

type FormatterOptions struct {
	Level      *logrus.Level
	Output     *OutputType
	Format     *FormatType
	ShowCaller bool
}

type CustomOptions struct {
	CustomFormat *CustomFormatter
}

type Config struct {
	FormatterOptions
	CustomOptions

	// Environment Mapping

	EnvLevel  string
	EnvOutput string
	EnvFormat string

	// Persistence

	Filename  string
	Namespace string // "LoggerName" is often called Namespace or Scope
}

func (c Config) setLevel(l *logrus.Logger) {
	// If user called WithLevel(), c.Level is not nil
	if c.Level != nil {
		l.SetLevel(*c.Level)
		return
	}

	// Fallback to Environment
	if env := os.Getenv(c.EnvLevel); env != "" {
		if p, err := logrus.ParseLevel(env); err == nil {
			l.SetLevel(p)
			return
		}
	}
	l.SetLevel(logrus.InfoLevel)
}

func (c Config) setOutput(l *logrus.Logger) {
	out := OutputConsole // Default
	if c.Output != nil {
		out = *c.Output
	} else if env := os.Getenv(c.EnvOutput); env != "" {
		out = parseOutputType(env)
	}

	switch out {
	case OutputFile:
		l.SetOutput(NewLumberjackLogger(c.Filename, DefaultLogFileConfig()))

	case OutputMulti:
		logFile := NewLumberjackLogger(c.Filename, DefaultLogFileConfig())

		// Create the multi-writer config using the resolved format
		mwConfig := MultiWriterWithFormattersConfig{
			format:       c.getFormat(),
			showCaller:   c.ShowCaller,
			customFormat: c.CustomFormat,
		}

		mw := NewMultiWriter(mwConfig)
		mw.AddWriter(os.Stdout, true, false) // Console gets colors
		mw.AddWriter(logFile, false, true)   // File gets timestamps, no colors

		l.AddHook(&CustomHook{mw: mw})
		l.SetOutput(io.Discard) // Hook handles writing

	default: // OutputConsole
		l.SetOutput(os.Stdout)
	}
}

func parseOutputType(env string) OutputType {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "file":
		return OutputFile
	case "multi":
		return OutputMulti
	case "console":
		return OutputConsole
	default:
		return OutputConsole
	}
}

func (c Config) setFormatter(l *logrus.Logger) {
	if c.CustomFormat != nil {
		l.SetFormatter(c.CustomFormat)
		l.SetReportCaller(c.CustomFormat.ShowCaller)
		return
	}

	fType := c.getFormat()
	l.SetReportCaller(c.ShowCaller)

	switch fType {
	case FormatJSON:
		l.SetFormatter(&logrus.JSONFormatter{})

	case FormatPlain:
		// If using Multi, the Hook handles formatting; don't set a global formatter
		isMulti := c.Output != nil && *c.Output == OutputMulti
		if !isMulti {
			useColors := c.Output != nil && *c.Output == OutputConsole
			l.SetFormatter(&CustomFormatter{
				UseColors:       useColors,
				ShowCaller:      c.ShowCaller,
				ShowTimestamp:   false,
				CallerLevel:     logrus.WarnLevel,
				UseRelativePath: true,
				BracketPadding:  15,
				ColorBrackets:   true,
			})
		}

	default: // FormatRaw
		l.SetFormatter(&logrus.TextFormatter{})
	}
}

// getFormat is a helper to resolve the format type from Struct -> Env -> Default
func (c Config) getFormat() FormatType {
	if c.Format != nil {
		return *c.Format
	}
	if env := os.Getenv(c.EnvFormat); env != "" {
		return parseFormatType(env)
	}
	return FormatRaw
}

func parseFormatType(env string) FormatType {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "multi":
		return FormatJSON
	case "console":
		return FormatPlain
	default:
		return FormatRaw
	}
}

func setup(l *logrus.Logger, cfg Config) {
	cfg.setLevel(l)
	cfg.setOutput(l)
	cfg.setFormatter(l)

	logInitComplete(l, cfg)
}

func logInitComplete(logger *logrus.Logger, cfg Config) {
	logger.Debugf("[Logger] %s initialized - Level: %s, Namespace: %s",
		cfg.Namespace,
		logger.GetLevel(),
		cfg.Namespace,
	)
}

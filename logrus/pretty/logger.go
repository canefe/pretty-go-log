package pretty

import (
	"github.com/sirupsen/logrus"
)

// Option is a function that modifies our Config
type Option func(*Config)

// New creates a logger by applying functional options to a default config
func New(opts ...Option) *logrus.Logger {
	// 1. Set defaults
	plain := FormatPlain
	console := OutputConsole
	cfg := &Config{
		FormatterOptions: FormatterOptions{
			Format:     &plain,
			Output:     &console,
			ShowCaller: true,
		},
		Namespace: "Main",
		EnvLevel:  "LOG_LEVEL",
		EnvOutput: "LOG_OUTPUT",
		EnvFormat: "LOG_FORMAT",
	}

	// 2. Apply user overrides
	for _, opt := range opts {
		opt(cfg)
	}

	l := logrus.New()
	setup(l, *cfg)
	return l
}

// --- Options functions ---

func WithLevel(l logrus.Level) Option {
	return func(c *Config) { c.Level = &l }
}

func WithOutput(o OutputType) Option {
	return func(c *Config) { c.Output = &o }
}

func WithFormat(f FormatType) Option {
	return func(c *Config) { c.Format = &f }
}

func WithCustomFormat(f CustomFormatter) Option {
	return func(c *Config) { c.CustomFormat = &f }
}

func WithNamespace(name string) Option {
	return func(c *Config) { c.Namespace = name }
}

func WithFile(path string) Option {
	return func(c *Config) { c.Filename = path }
}

func WithoutCaller() Option {
	return func(c *Config) { c.ShowCaller = false }
}

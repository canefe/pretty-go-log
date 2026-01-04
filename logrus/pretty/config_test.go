package pretty

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestConfig_setLevel_Explicit(t *testing.T) {
	logger := logrus.New()
	level := logrus.WarnLevel
	cfg := Config{FormatterOptions: FormatterOptions{Level: &level}}

	cfg.setLevel(logger)

	if logger.GetLevel() != logrus.WarnLevel {
		t.Errorf("Expected warn level, got %v", logger.GetLevel())
	}
}

func TestConfig_setLevel_FromEnv(t *testing.T) {
	t.Setenv("TEST_LOG_LEVEL_ENV", "debug")

	logger := logrus.New()
	cfg := Config{EnvLevel: "TEST_LOG_LEVEL_ENV"}
	cfg.setLevel(logger)

	if logger.GetLevel() != logrus.DebugLevel {
		t.Errorf("Expected debug level from env, got %v", logger.GetLevel())
	}
}

func TestConfig_setLevel_Default(t *testing.T) {
	logger := logrus.New()
	cfg := Config{EnvLevel: "TEST_LOG_LEVEL_ENV"}
	cfg.setLevel(logger)

	if logger.GetLevel() != logrus.InfoLevel {
		t.Errorf("Expected default InfoLevel, got %v", logger.GetLevel())
	}
}

func TestConfig_setOutput_Console(t *testing.T) {
	logger := logrus.New()
	output := OutputConsole
	cfg := Config{FormatterOptions: FormatterOptions{Output: &output}}

	cfg.setOutput(logger)

	if logger.Out != os.Stdout {
		t.Error("Expected output to be stdout for console output")
	}
}

func TestConfig_setOutput_File(t *testing.T) {
	logger := logrus.New()
	output := OutputFile
	cfg := Config{
		FormatterOptions: FormatterOptions{Output: &output},
		Filename:         "test.log",
	}

	cfg.setOutput(logger)

	if _, ok := logger.Out.(*lumberjack.Logger); !ok {
		t.Error("Expected lumberjack logger for file output")
	}
}

func TestConfig_setOutput_Multi(t *testing.T) {
	logger := logrus.New()
	output := OutputMulti
	format := FormatPlain
	cfg := Config{
		FormatterOptions: FormatterOptions{Output: &output, Format: &format},
		Filename:         "test.log",
	}

	cfg.setOutput(logger)

	if logger.Out != io.Discard {
		t.Error("Expected output to be discarded for multi output")
	}

	foundHook := false
	for _, hooks := range logger.Hooks {
		for _, hook := range hooks {
			if _, ok := hook.(*CustomHook); ok {
				foundHook = true
				break
			}
		}
	}
	if !foundHook {
		t.Error("Expected CustomHook to be registered for multi output")
	}
}

func TestParseOutputType(t *testing.T) {
	if got := parseOutputType("file"); got != OutputFile {
		t.Errorf("Expected OutputFile, got %v", got)
	}
	if got := parseOutputType("multi"); got != OutputMulti {
		t.Errorf("Expected OutputMulti, got %v", got)
	}
	if got := parseOutputType("console"); got != OutputConsole {
		t.Errorf("Expected OutputConsole, got %v", got)
	}
	if got := parseOutputType("unknown"); got != OutputConsole {
		t.Errorf("Expected OutputConsole for unknown output, got %v", got)
	}
}

func TestConfig_setFormatter_JSON(t *testing.T) {
	logger := logrus.New()
	format := FormatJSON
	cfg := Config{FormatterOptions: FormatterOptions{Format: &format}}

	cfg.setFormatter(logger)

	if _, ok := logger.Formatter.(*logrus.JSONFormatter); !ok {
		t.Error("Expected JSON formatter")
	}
}

func TestConfig_setFormatter_Plain(t *testing.T) {
	logger := logrus.New()
	format := FormatPlain
	output := OutputConsole
	cfg := Config{FormatterOptions: FormatterOptions{Format: &format, Output: &output}}

	cfg.setFormatter(logger)

	if _, ok := logger.Formatter.(*CustomFormatter); !ok {
		t.Error("Expected CustomFormatter for plain format")
	}
}

func TestConfig_setFormatter_Raw(t *testing.T) {
	logger := logrus.New()
	format := FormatRaw
	cfg := Config{FormatterOptions: FormatterOptions{Format: &format}}

	cfg.setFormatter(logger)

	if _, ok := logger.Formatter.(*logrus.TextFormatter); !ok {
		t.Error("Expected TextFormatter for raw format")
	}
}

func TestConfig_setFormatter_PlainMultiOutput(t *testing.T) {
	logger := logrus.New()
	format := FormatPlain
	output := OutputMulti
	cfg := Config{FormatterOptions: FormatterOptions{Format: &format, Output: &output}}

	cfg.setFormatter(logger)

	if _, ok := logger.Formatter.(*logrus.TextFormatter); !ok {
		t.Error("Expected TextFormatter when using plain format with multi output")
	}
}

func TestConfig_getFormat_FromEnv(t *testing.T) {
	t.Setenv("TEST_FORMAT_ENV", "console")

	cfg := Config{EnvFormat: "TEST_FORMAT_ENV"}
	if got := cfg.getFormat(); got != FormatPlain {
		t.Errorf("Expected FormatPlain from env, got %v", got)
	}
}

func TestLogInitComplete(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.DebugLevel)

	cfg := Config{Namespace: "TestLogger"}
	logInitComplete(logger, cfg)

	if !bytes.Contains(buf.Bytes(), []byte("TestLogger")) {
		t.Error("Expected logger namespace in initialization info")
	}
}

package pretty

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestNew_Defaults(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("LOG_OUTPUT", "")
	t.Setenv("LOG_FORMAT", "")

	logger := New()

	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	if logger.GetLevel() != logrus.InfoLevel {
		t.Errorf("Expected default level InfoLevel, got %v", logger.GetLevel())
	}

	if !logger.ReportCaller {
		t.Error("Expected ReportCaller to be true by default")
	}

	if _, ok := logger.Formatter.(*CustomFormatter); !ok {
		t.Error("Expected CustomFormatter by default")
	}

	if logger.Out != os.Stdout {
		t.Error("Expected output to be stdout by default")
	}
}

func TestNew_WithLevel(t *testing.T) {
	logger := New(WithLevel(logrus.DebugLevel))

	if logger.GetLevel() != logrus.DebugLevel {
		t.Errorf("Expected debug level, got %v", logger.GetLevel())
	}
}

func TestNew_WithFormatJSON(t *testing.T) {
	logger := New(
		WithFormat(FormatJSON),
		WithOutput(OutputConsole),
	)

	if _, ok := logger.Formatter.(*logrus.JSONFormatter); !ok {
		t.Error("Expected JSONFormatter for json format")
	}
}

func TestNew_WithFormatPlain(t *testing.T) {
	logger := New(
		WithFormat(FormatPlain),
		WithOutput(OutputConsole),
	)

	if _, ok := logger.Formatter.(*CustomFormatter); !ok {
		t.Error("Expected CustomFormatter for plain format")
	}
}

func TestNew_WithOutputFile(t *testing.T) {
	logger := New(
		WithLevel(logrus.InfoLevel),
		WithOutput(OutputFile),
		WithFormat(FormatPlain),
		WithFile("test.log"),
	)

	if _, ok := logger.Out.(*lumberjack.Logger); !ok {
		t.Error("Expected lumberjack logger for file output")
	}
}

func TestNew_WithOutputMulti(t *testing.T) {
	logger := New(
		WithLevel(logrus.InfoLevel),
		WithOutput(OutputMulti),
		WithFormat(FormatPlain),
		WithFile("test.log"),
	)

	if logger.Out != io.Discard {
		t.Error("Expected output to be discarded when using multi output")
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

func TestNew_WithoutCaller(t *testing.T) {
	logger := New(
		WithFormat(FormatPlain),
		WithOutput(OutputConsole),
		WithoutCaller(),
	)

	if logger.ReportCaller {
		t.Error("Expected ReportCaller to be false when WithoutCaller is used")
	}
}

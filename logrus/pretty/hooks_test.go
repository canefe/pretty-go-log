package pretty

import (
	"bytes"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewMultiWriter(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatPlain,
		showCaller: true,
	}

	mw := NewMultiWriter(config)

	if mw == nil {
		t.Fatal("NewMultiWriter returned nil")
	}

	if len(mw.pairs) != 0 {
		t.Errorf("Expected 0 writer pairs initially, got %d", len(mw.pairs))
	}

	if mw.cfg.format != FormatPlain {
		t.Errorf("Expected format FormatPlain, got %v", mw.cfg.format)
	}

	if !mw.cfg.showCaller {
		t.Error("Expected showCaller to be true")
	}
}

func TestMultiWriter_AddWriter(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatPlain,
		showCaller: false,
	}

	mw := NewMultiWriter(config)
	var buf bytes.Buffer

	mw.AddWriter(&buf, true, false)

	if len(mw.pairs) != 1 {
		t.Errorf("Expected 1 writer pair after AddWriter, got %d", len(mw.pairs))
	}

	// Verify formatter is set
	if mw.pairs[0].f == nil {
		t.Error("Expected formatter to be set")
	}

	if mw.pairs[0].w != &buf {
		t.Error("Expected writer to be the buffer")
	}
}

func TestMultiWriter_AddMultipleWriters(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatPlain,
		showCaller: true,
	}

	mw := NewMultiWriter(config)
	var buf1, buf2 bytes.Buffer

	mw.AddWriter(&buf1, true, false)
	mw.AddWriter(&buf2, false, true)

	if len(mw.pairs) != 2 {
		t.Errorf("Expected 2 writer pairs, got %d", len(mw.pairs))
	}
}

func TestMultiWriter_WriteEntry(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatPlain,
		showCaller: false,
	}

	mw := NewMultiWriter(config)
	var buf bytes.Buffer

	mw.AddWriter(&buf, false, false)

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "test message"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	err := mw.WriteEntry(entry)
	if err != nil {
		t.Errorf("WriteEntry failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Expected output to be written to buffer")
	}
}

func TestMultiWriter_JSONFormat(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatJSON,
		showCaller: false,
	}

	mw := NewMultiWriter(config)
	var buf bytes.Buffer

	mw.AddWriter(&buf, false, false)

	// Verify JSON formatter is used
	if _, ok := mw.pairs[0].f.(*logrus.JSONFormatter); !ok {
		t.Error("Expected JSONFormatter for json format")
	}
}

func TestMultiWriter_PlainFormat(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatPlain,
		showCaller: true,
	}

	mw := NewMultiWriter(config)
	var buf bytes.Buffer

	mw.AddWriter(&buf, true, false)

	// Verify CustomFormatter is used
	if _, ok := mw.pairs[0].f.(*CustomFormatter); !ok {
		t.Error("Expected CustomFormatter for plain format")
	}
}

func TestMultiWriter_DefaultFormat(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatRaw,
		showCaller: false,
	}

	mw := NewMultiWriter(config)
	var buf bytes.Buffer

	mw.AddWriter(&buf, false, false)

	// Verify TextFormatter is used for raw format
	if _, ok := mw.pairs[0].f.(*logrus.TextFormatter); !ok {
		t.Error("Expected TextFormatter for raw format")
	}
}

func TestMultiWriter_CustomFormat(t *testing.T) {
	format := CustomFormatter{
		BracketPadding: 42,
	}
	config := MultiWriterWithFormattersConfig{
		format:       FormatPlain,
		showCaller:   true,
		customFormat: &format,
	}

	mw := NewMultiWriter(config)
	var buf bytes.Buffer

	mw.AddWriter(&buf, false, false)

	cf, ok := mw.pairs[0].f.(*CustomFormatter)
	if !ok {
		t.Fatal("Expected CustomFormatter for custom format")
	}

	if cf.BracketPadding != 42 {
		t.Errorf("Expected custom BracketPadding 42, got %d", cf.BracketPadding)
	}
}

func TestCustomHook_Levels(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatPlain,
		showCaller: false,
	}

	mw := NewMultiWriter(config)
	hook := &CustomHook{mw: mw}

	levels := hook.Levels()

	if len(levels) != len(logrus.AllLevels) {
		t.Errorf("Expected %d levels, got %d", len(logrus.AllLevels), len(levels))
	}
}

func TestCustomHook_Fire(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatPlain,
		showCaller: false,
	}

	mw := NewMultiWriter(config)
	var buf bytes.Buffer
	mw.AddWriter(&buf, false, false)

	hook := &CustomHook{mw: mw}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "hook test"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	err := hook.Fire(entry)
	if err != nil {
		t.Errorf("Fire failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Expected output to be written via hook")
	}
}

func TestCustomHook_Integration(t *testing.T) {
	config := MultiWriterWithFormattersConfig{
		format:     FormatPlain,
		showCaller: false,
	}

	mw := NewMultiWriter(config)
	var buf bytes.Buffer
	mw.AddWriter(&buf, false, false)

	hook := &CustomHook{mw: mw}

	logger := logrus.New()
	logger.AddHook(hook)
	logger.SetOutput(&bytes.Buffer{}) // Discard default output

	logger.Info("test message via hook")

	if buf.Len() == 0 {
		t.Error("Expected output to be written via hook integration")
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected non-empty output")
	}
}

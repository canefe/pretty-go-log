package pretty

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestFormatter_Format(t *testing.T) {
	f := &CustomFormatter{CallerLevel: logrus.WarnLevel}
	l := logrus.New()
	entry := l.WithField("user", "canefe")
	entry.Message = "test message"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	if len(b) == 0 {
		t.Error("Resulting byte slice is empty")
	}
}

func TestFormatter_Brackets(t *testing.T) {
	f := &CustomFormatter{UseColors: false, CallerLevel: logrus.WarnLevel}
	l := logrus.New()
	entry := l.WithField("user", "canefe")
	entry.Message = "[TestTag] test message with bracketed tag"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Output missing brackets: %s", output)
	}
}

func TestFormatter_WithColors(t *testing.T) {
	f := &CustomFormatter{UseColors: true, CallerLevel: logrus.WarnLevel}
	l := logrus.New()
	entry := l.WithField("user", "canefe")
	entry.Message = "colored message"
	entry.Level = logrus.WarnLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should contain ANSI color codes when UseColors is true
	if !strings.Contains(output, "\033[") {
		t.Error("Output missing color codes when UseColors=true")
	}
}

func TestFormatter_WithoutColors(t *testing.T) {
	f := &CustomFormatter{UseColors: false, CallerLevel: logrus.WarnLevel, UseRelativePath: true}
	l := logrus.New()
	entry := l.WithField("user", "canefe")
	entry.Message = "plain message"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should not contain ANSI color codes when UseColors is false
	if strings.Contains(output, "\033[") {
		t.Error("Output contains color codes when UseColors=false")
	}
}

func TestFormatter_WithTimestamp(t *testing.T) {
	f := &CustomFormatter{ShowTimestamp: true, CallerLevel: logrus.WarnLevel}
	l := logrus.New()
	entry := l.WithField("user", "canefe")
	entry.Message = "timestamped message"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should contain timestamp in brackets
	if !strings.Contains(output, "[20") {
		t.Error("Output missing timestamp when ShowTimestamp=true")
	}
}

func TestFormatter_WithoutTimestamp(t *testing.T) {
	f := &CustomFormatter{ShowTimestamp: false, CallerLevel: logrus.WarnLevel}
	l := logrus.New()
	entry := l.WithField("user", "canefe")
	entry.Message = "no timestamp message"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should not start with a timestamp bracket
	if strings.HasPrefix(output, "[20") {
		t.Error("Output contains timestamp when ShowTimestamp=false")
	}
}

func TestFormatter_WithCaller(t *testing.T) {
	f := &CustomFormatter{ShowCaller: true, CallerLevel: logrus.WarnLevel}
	l := logrus.New()
	l.SetReportCaller(true)

	// Log an actual error to get real caller info
	var buf strings.Builder
	l.SetOutput(&buf)
	l.SetFormatter(f)
	l.Error("error with caller")

	output := buf.String()

	// Error and above should show caller info when ShowCaller=true
	if !strings.Contains(output, "at (") {
		t.Error("Output missing caller info for ERROR level when ShowCaller=true")
	}
}

func TestFormatter_LevelFormatting(t *testing.T) {
	tests := []struct {
		level        logrus.Level
		expectedText string
	}{
		{logrus.DebugLevel, "DEBUG"},
		{logrus.InfoLevel, "INFO"},
		{logrus.WarnLevel, "WARN"},
		{logrus.ErrorLevel, "ERROR"},
		{logrus.FatalLevel, "FATAL"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedText, func(t *testing.T) {
			f := &CustomFormatter{UseColors: false, CallerLevel: logrus.WarnLevel}
			l := logrus.New()
			entry := l.WithField("test", "value")
			entry.Message = "test message"
			entry.Level = tt.level
			entry.Time = time.Now()

			b, err := f.Format(entry)
			if err != nil {
				t.Fatalf("Failed to format: %v", err)
			}

			output := string(b)
			if !strings.Contains(output, tt.expectedText) {
				t.Errorf("Expected output to contain %s, got: %s", tt.expectedText, output)
			}
		})
	}
}

func TestFormatter_BracketPadding(t *testing.T) {
	f := &CustomFormatter{UseColors: false, CallerLevel: logrus.WarnLevel}
	l := logrus.New()

	tests := []struct {
		message string
		name    string
	}{
		{"[Short] message", "short tag"},
		{"[VeryLongTagName] message", "long tag"},
		{"[MediumTag] message", "medium tag"},
		{"no brackets here", "no brackets"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := l.WithField("test", "value")
			entry.Message = tt.message
			entry.Level = logrus.InfoLevel
			entry.Time = time.Now()

			b, err := f.Format(entry)
			if err != nil {
				t.Fatalf("Failed to format: %v", err)
			}

			if len(b) == 0 {
				t.Error("Resulting byte slice is empty")
			}
		})
	}
}

func TestFormatCommon(t *testing.T) {
	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "test message"
	entry.Level = logrus.WarnLevel
	entry.Time = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create formatter with colors enabled
	f := &CustomFormatter{UseColors: true, ShowTimestamp: true}
	timestamp, level, message, colorCode, resetCode := f.formatCommon(entry)

	if timestamp != "2024-01-01 12:00:00" {
		t.Errorf("Expected timestamp '2024-01-01 12:00:00', got '%s'", timestamp)
	}

	if level != "WARN" {
		t.Errorf("Expected level 'WARN', got '%s'", level)
	}

	if message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", message)
	}

	if colorCode == "" {
		t.Error("Expected color code when useColors=true")
	}

	if resetCode != "\033[0m" {
		t.Errorf("Expected reset code '\\033[0m', got '%s'", resetCode)
	}
}

func TestFormatCommon_NoColors(t *testing.T) {
	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "test message"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	// Create formatter with colors disabled
	f := &CustomFormatter{UseColors: false, ShowTimestamp: false}
	_, _, _, colorCode, resetCode := f.formatCommon(entry)

	if colorCode != "" {
		t.Error("Expected no color code when useColors=false")
	}

	if resetCode != "" {
		t.Error("Expected no reset code when useColors=false")
	}
}

func TestFormatter_CallerLevelThreshold(t *testing.T) {
	tests := []struct {
		name        string
		callerLevel logrus.Level
		entryLevel  logrus.Level
		shouldShow  bool
	}{
		{"Show for Error when CallerLevel=Warn", logrus.WarnLevel, logrus.ErrorLevel, true},
		{"Show for Warn when CallerLevel=Warn", logrus.WarnLevel, logrus.WarnLevel, true},
		{"Hide for Info when CallerLevel=Warn", logrus.WarnLevel, logrus.InfoLevel, false},
		{"Hide for Debug when CallerLevel=Warn", logrus.WarnLevel, logrus.DebugLevel, false},
		{"Show for Info when CallerLevel=Info", logrus.InfoLevel, logrus.InfoLevel, true},
		{"Show for Error when CallerLevel=Error", logrus.ErrorLevel, logrus.ErrorLevel, true},
		{"Hide for Warn when CallerLevel=Error", logrus.ErrorLevel, logrus.WarnLevel, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &CustomFormatter{
				ShowCaller:  true,
				CallerLevel: tt.callerLevel,
				UseColors:   false,
			}

			l := logrus.New()
			l.SetReportCaller(true)
			var buf strings.Builder
			l.SetOutput(&buf)
			l.SetFormatter(f)

			// Log at the specified level
			switch tt.entryLevel {
			case logrus.DebugLevel:
				l.Debug("test message")
			case logrus.InfoLevel:
				l.Info("test message")
			case logrus.WarnLevel:
				l.Warn("test message")
			case logrus.ErrorLevel:
				l.Error("test message")
			}

			output := buf.String()
			containsCaller := strings.Contains(output, "at (")

			if tt.shouldShow && !containsCaller {
				t.Errorf("Expected caller info to be shown for %s at CallerLevel=%s", tt.entryLevel, tt.callerLevel)
			}

			if !tt.shouldShow && containsCaller {
				t.Errorf("Expected caller info to be hidden for %s at CallerLevel=%s", tt.entryLevel, tt.callerLevel)
			}
		})
	}
}

func TestNewCustomFormatter_Defaults(t *testing.T) {
	f := NewCustomFormatter()

	if !f.UseColors {
		t.Error("Expected UseColors to be true by default")
	}

	if !f.ShowCaller {
		t.Error("Expected ShowCaller to be true by default")
	}

	if f.ShowTimestamp {
		t.Error("Expected ShowTimestamp to be false by default")
	}

	if f.CallerLevel != logrus.WarnLevel {
		t.Errorf("Expected CallerLevel to be WarnLevel by default, got %v", f.CallerLevel)
	}
}

func TestNewCustomFormatter_WithOptions(t *testing.T) {
	f := NewCustomFormatter(
		WithColors(false),
		WithTimestamp(true),
		WithCaller(true, logrus.InfoLevel),
	)

	if f.UseColors {
		t.Error("Expected UseColors to be false")
	}

	if !f.ShowTimestamp {
		t.Error("Expected ShowTimestamp to be true")
	}

	if !f.ShowCaller {
		t.Error("Expected ShowCaller to be true")
	}

	if f.CallerLevel != logrus.InfoLevel {
		t.Errorf("Expected CallerLevel to be InfoLevel, got %v", f.CallerLevel)
	}
}

func TestFormatterOptions_Individual(t *testing.T) {
	// Test WithColors
	f1 := NewCustomFormatter(WithColors(false))
	if f1.UseColors {
		t.Error("WithColors(false) did not work")
	}

	// Test WithTimestamp
	f2 := NewCustomFormatter(WithTimestamp(true))
	if !f2.ShowTimestamp {
		t.Error("WithTimestamp(true) did not work")
	}

	// Test WithCaller
	f3 := NewCustomFormatter(WithCaller(false, logrus.ErrorLevel))
	if f3.ShowCaller {
		t.Error("WithCaller(false, ...) did not work")
	}
	if f3.CallerLevel != logrus.ErrorLevel {
		t.Errorf("WithCaller(..., ErrorLevel) did not work, got %v", f3.CallerLevel)
	}
}

func TestFormatter_WithFields(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
	}

	l := logrus.New()
	entry := l.WithFields(map[string]interface{}{
		"user":     "alice",
		"duration": "150ms",
	})
	entry.Message = "test message"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should contain fields (checking for key= pattern since value might have color codes)
	if !strings.Contains(output, "duration=") {
		t.Errorf("Output missing duration field, got: %s", output)
	}

	if !strings.Contains(output, "user=") {
		t.Errorf("Output missing user field, got: %s", output)
	}

	// Verify fields are in sorted order (duration comes before user)
	durationIdx := strings.Index(output, "duration=")
	userIdx := strings.Index(output, "user=")
	if durationIdx >= userIdx {
		t.Error("Fields not in sorted order: expected duration before user")
	}
}

func TestFormatter_WithFieldsColored(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       true,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
	}

	l := logrus.New()
	entry := l.WithField("request_id", "abc123")
	entry.Message = "API request"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should contain the field (just check for key= since value has colors)
	if !strings.Contains(output, "request_id=") {
		t.Errorf("Output missing request_id field, got: %s", output)
	}

	// When colors are enabled, fields should be in dim color
	if !strings.Contains(output, ColorVeryDimGray) {
		t.Error("Fields missing dim color when UseColors=true")
	}
}

func TestFormatter_NoFieldsEmptyOutput(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
	}

	l := logrus.New()
	entry := l.WithTime(time.Now())
	entry.Message = "simple message"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should not have extra spaces or equals signs
	if strings.Contains(output, "=") {
		t.Error("Output contains '=' when there are no fields")
	}

	// Should contain the message
	if !strings.Contains(output, "simple message") {
		t.Error("Output missing message")
	}
}

func TestFormatter_MultipleFieldsOrdering(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
	}

	l := logrus.New()
	// Add fields in non-alphabetical order
	entry := l.WithFields(map[string]interface{}{
		"zebra":  "last",
		"apple":  "first",
		"middle": "middle",
		"banana": "second",
	})
	entry.Message = "sorting test"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Find positions of each field
	appleIdx := strings.Index(output, "apple=")
	bananaIdx := strings.Index(output, "banana=")
	middleIdx := strings.Index(output, "middle=")
	zebraIdx := strings.Index(output, "zebra=")

	// Verify alphabetical ordering
	if !(appleIdx < bananaIdx && bananaIdx < middleIdx && middleIdx < zebraIdx) {
		t.Errorf("Fields not in alphabetical order. Positions: apple=%d, banana=%d, middle=%d, zebra=%d",
			appleIdx, bananaIdx, middleIdx, zebraIdx)
	}
}

func TestFormatter_FieldWithBracketedMessage(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
	}

	l := logrus.New()
	entry := l.WithField("endpoint", "/api/users")
	entry.Message = "[Request] GET /api/users"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should contain both the bracketed tag and the field
	if !strings.Contains(output, "[Request]") {
		t.Error("Output missing bracketed tag")
	}

	if !strings.Contains(output, "endpoint=") {
		t.Errorf("Output missing field, got: %s", output)
	}

	// Message and field should be separated
	msgIdx := strings.Index(output, "GET /api/users")
	fieldIdx := strings.Index(output, "endpoint=")
	if fieldIdx <= msgIdx {
		t.Error("Field should appear after the message")
	}
}

func TestFormatter_RelativePath(t *testing.T) {
	f := &CustomFormatter{
		ShowCaller:      true,
		CallerLevel:     logrus.WarnLevel,
		UseColors:       false,
		UseRelativePath: true,
	}

	l := logrus.New()
	l.SetReportCaller(true)
	var buf strings.Builder
	l.SetOutput(&buf)
	l.SetFormatter(f)
	l.Error("error with relative path")

	output := buf.String()

	// Should contain "at (" but NOT absolute path
	if !strings.Contains(output, "at (") {
		t.Error("Output missing caller info")
	}

	// Relative path should not start with "/"
	// Extract the path between "at (" and ":"
	startIdx := strings.Index(output, "at (")
	if startIdx != -1 {
		pathPart := output[startIdx+4:]
		colonIdx := strings.Index(pathPart, ":")
		if colonIdx != -1 {
			path := pathPart[:colonIdx]
			if strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "../") {
				t.Errorf("Expected relative path, got absolute path: %s", path)
			}
		}
	}
}

func TestFormatter_AbsolutePath(t *testing.T) {
	f := &CustomFormatter{
		ShowCaller:      true,
		CallerLevel:     logrus.WarnLevel,
		UseColors:       false,
		UseRelativePath: false,
	}

	l := logrus.New()
	l.SetReportCaller(true)
	var buf strings.Builder
	l.SetOutput(&buf)
	l.SetFormatter(f)
	l.Error("error with absolute path")

	output := buf.String()

	// Should contain "at (" with absolute path
	if !strings.Contains(output, "at (") {
		t.Error("Output missing caller info")
	}

	// Absolute path should start with "/"
	startIdx := strings.Index(output, "at (")
	if startIdx != -1 {
		pathPart := output[startIdx+4:]
		colonIdx := strings.Index(pathPart, ":")
		if colonIdx != -1 {
			path := pathPart[:colonIdx]
			if !strings.HasPrefix(path, "/") {
				t.Errorf("Expected absolute path starting with '/', got: %s", path)
			}
		}
	}
}

func TestWithRelativePath_Option(t *testing.T) {
	// Test enabling relative paths
	f1 := NewCustomFormatter(WithRelativePath(true))
	if !f1.UseRelativePath {
		t.Error("WithRelativePath(true) did not work")
	}

	// Test disabling relative paths (use absolute)
	f2 := NewCustomFormatter(WithRelativePath(false))
	if f2.UseRelativePath {
		t.Error("WithRelativePath(false) did not work")
	}
}

func TestNewCustomFormatter_RelativePathDefault(t *testing.T) {
	f := NewCustomFormatter()

	// Default should be relative paths
	if !f.UseRelativePath {
		t.Error("Expected UseRelativePath to be true by default")
	}
}

func TestWithBracketPadding_Option(t *testing.T) {
	// Test custom padding
	f1 := NewCustomFormatter(WithBracketPadding(20))
	if f1.BracketPadding != 20 {
		t.Errorf("WithBracketPadding(20) did not work, got %d", f1.BracketPadding)
	}

	// Test zero padding
	f2 := NewCustomFormatter(WithBracketPadding(0))
	if f2.BracketPadding != 0 {
		t.Errorf("WithBracketPadding(0) did not work, got %d", f2.BracketPadding)
	}

	// Test negative padding (should be clamped to 0)
	f3 := NewCustomFormatter(WithBracketPadding(-5))
	if f3.BracketPadding != 0 {
		t.Errorf("WithBracketPadding(-5) should clamp to 0, got %d", f3.BracketPadding)
	}
}

func TestWithColorBrackets_Option(t *testing.T) {
	// Test enabling bracket coloring
	f1 := NewCustomFormatter(WithColorBrackets(true))
	if !f1.ColorBrackets {
		t.Error("WithColorBrackets(true) did not work")
	}

	// Test disabling bracket coloring
	f2 := NewCustomFormatter(WithColorBrackets(false))
	if f2.ColorBrackets {
		t.Error("WithColorBrackets(false) did not work")
	}
}

func TestNewCustomFormatter_BracketDefaults(t *testing.T) {
	f := NewCustomFormatter()

	if f.BracketPadding != 15 {
		t.Errorf("Expected BracketPadding to be 15 by default, got %d", f.BracketPadding)
	}

	if !f.ColorBrackets {
		t.Error("Expected ColorBrackets to be true by default")
	}
}

func TestFormatter_BracketPaddingConfig(t *testing.T) {
	tests := []struct {
		name      string
		padding   int
		message   string
		expectPad bool
	}{
		{"15 char padding with short tag", 15, "[API] message", true},
		{"10 char padding with short tag", 10, "[API] message", true},
		{"5 char padding with long tag", 5, "[VeryLongTag] message", false},
		{"0 padding (disabled)", 0, "[API] message", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &CustomFormatter{
				UseColors:       false,
				CallerLevel:     logrus.WarnLevel,
				UseRelativePath: true,
				BracketPadding:  tt.padding,
				ColorBrackets:   false,
			}

			l := logrus.New()
			entry := l.WithField("test", "value")
			entry.Message = tt.message
			entry.Level = logrus.InfoLevel
			entry.Time = time.Now()

			b, err := f.Format(entry)
			if err != nil {
				t.Fatalf("Failed to format: %v", err)
			}

			output := string(b)

			// Should contain the bracketed tag
			if !strings.Contains(output, "[") {
				t.Error("Output missing bracket")
			}
		})
	}
}

func TestFormatter_ColorBrackets_Enabled(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       true,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   true,
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[Server] Starting server"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should contain yellow color code for brackets when ColorBrackets=true
	if !strings.Contains(output, ColorYellow) {
		t.Error("Output missing bracket color when ColorBrackets=true")
	}
}

func TestFormatter_ColorBrackets_Disabled(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       true,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   false,
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[Server] Starting server"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Extract the message part (after level indicator)
	// Should not have ColorYellow on the bracket when ColorBrackets=false
	// But level colors will still be present, so check specifically for bracket region
	tagStart := strings.Index(output, "[Server]")
	if tagStart == -1 {
		t.Fatal("Could not find [Server] tag in output")
	}

	// Check if there's yellow color right before the bracket
	beforeBracket := output[max(0, tagStart-20):tagStart]
	if strings.Contains(beforeBracket, ColorYellow) {
		t.Error("Brackets should not be colored when ColorBrackets=false")
	}
}

func TestFormatter_BracketPaddingZero(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  0,
		ColorBrackets:   false,
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[API] message"
	entry.Level = logrus.InfoLevel
	entry.Time = time.Now()

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format: %v", err)
	}

	output := string(b)

	// Should still contain the message and bracket
	if !strings.Contains(output, "[API]") {
		t.Error("Output missing bracket tag")
	}

	if !strings.Contains(output, "message") {
		t.Error("Output missing message")
	}
}

func TestWithCenterBrackets_Option(t *testing.T) {
	f := NewCustomFormatter(WithCenterBrackets(true, "="))
	if f.TagStyle != StyleCenter {
		t.Error("Expected TagStyle to be StyleCenter")
	}
	if f.PaddingChar != "=" {
		t.Errorf("Expected PaddingChar '=', got '%s'", f.PaddingChar)
	}

	// Test with different char
	f2 := NewCustomFormatter(WithCenterBrackets(true, "-"))
	if f2.PaddingChar != "-" {
		t.Errorf("Expected PaddingChar '-', got '%s'", f2.PaddingChar)
	}

	// Test with empty char (should default to "•")
	f3 := NewCustomFormatter(WithCenterBrackets(true, ""))
	if f3.PaddingChar != "•" {
		t.Errorf("Expected default PaddingChar '•', got '%s'", f3.PaddingChar)
	}
}

func TestFormatter_CenterBrackets_ShortTag(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   false,
		TagStyle:        StyleCenter,
		PaddingChar:     "=",
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[Init] Starting application"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	output := string(b)
	// "Init" is 4 chars, maxPadding is 15, so:
	// availableSpace = 15 - 2 = 13
	// sidePad = (13 - 4 - 2) / 2 = 3.5 -> 3
	// Result should be: [=== Init ===] with extra = on right for odd spacing
	if !strings.Contains(output, "[=== Init ===]") && !strings.Contains(output, "[=== Init ====]") {
		t.Errorf("Expected centered tag with '=' padding, got: %s", output)
	}
}

func TestFormatter_CenterBrackets_DifferentChars(t *testing.T) {
	tests := []struct {
		name        string
		paddingChar string
		tag         string
		expected    string
	}{
		{"equals", "=", "Init", "=== Init ==="},
		{"dash", "-", "Init", "--- Init ---"},
		{"dot", "·", "Init", "··· Init ···"},
		{"bullet", "•", "Init", "••• Init •••"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &CustomFormatter{
				UseColors:       false,
				CallerLevel:     logrus.WarnLevel,
				UseRelativePath: true,
				BracketPadding:  15,
				ColorBrackets:   false,
				TagStyle:        StyleCenter,
				PaddingChar:     tt.paddingChar,
			}

			l := logrus.New()
			entry := l.WithField("test", "value")
			entry.Message = fmt.Sprintf("[%s] Starting", tt.tag)
			entry.Level = logrus.InfoLevel

			b, err := f.Format(entry)
			if err != nil {
				t.Fatalf("Format error: %v", err)
			}

			output := string(b)
			// Check if the expected pattern is in the output (allow for odd spacing)
			if !strings.Contains(output, "["+tt.expected) && !strings.Contains(output, "["+tt.expected+tt.paddingChar) {
				t.Errorf("Expected pattern '[%s' in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestFormatter_CenterBrackets_VariousLengths(t *testing.T) {
	tests := []struct {
		tag    string
		minPad int // Minimum padding chars on each side
		maxPad int // Maximum padding chars on each side
	}{
		{"DB", 4, 5},       // Very short: 2 chars
		{"Init", 3, 4},     // Short: 4 chars
		{"Server", 2, 3},   // Medium: 6 chars
		{"Database", 1, 2}, // Longer: 8 chars
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			f := &CustomFormatter{
				UseColors:       false,
				CallerLevel:     logrus.WarnLevel,
				UseRelativePath: true,
				BracketPadding:  15,
				ColorBrackets:   false,
				TagStyle:        StyleCenter,
				PaddingChar:     "=",
			}

			l := logrus.New()
			entry := l.WithField("test", "value")
			entry.Message = fmt.Sprintf("[%s] Message", tt.tag)
			entry.Level = logrus.InfoLevel

			b, err := f.Format(entry)
			if err != nil {
				t.Fatalf("Format error: %v", err)
			}

			output := string(b)
			// Verify the tag is present and has padding
			if !strings.Contains(output, "[") || !strings.Contains(output, tt.tag) {
				t.Errorf("Expected tag '[%s' in output, got: %s", tt.tag, output)
			}

			// Verify there's at least some padding
			if !strings.Contains(output, "= "+tt.tag) && !strings.Contains(output, "==") {
				t.Errorf("Expected padding around tag, got: %s", output)
			}
		})
	}
}

func TestFormatter_CenterBrackets_Disabled(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   false,
		TagStyle:        StyleDefault, // Disabled
		PaddingChar:     "=",
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[Init] Starting"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	output := string(b)
	// Should NOT have centered padding when disabled
	if strings.Contains(output, "===") {
		t.Errorf("Should not have centered padding when CenterBrackets=false, got: %s", output)
	}
	// Should have plain bracket
	if !strings.Contains(output, "[Init]") {
		t.Errorf("Expected plain '[Init]', got: %s", output)
	}
}

func TestFormatter_CenterBrackets_WithColors(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       true,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   true,
		TagStyle:        StyleCenter,
		PaddingChar:     "=",
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[Init] Starting"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	output := string(b)
	// Should have both centered padding AND yellow color
	if !strings.Contains(output, ColorYellow) {
		t.Error("Expected yellow color when UseColors=true and ColorBrackets=true")
	}
	clean := stripANSI(output)
	if !strings.Contains(clean, "=== Init ===") && !strings.Contains(clean, "=== Init ====") {
		t.Errorf("Expected centered tag, got: %s", output)
	}
	if !strings.Contains(output, ColorVeryDimGray) {
		t.Error("Expected dim gray padding color")
	}
}

func TestFormatter_CenterBrackets_LongTag(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   false,
		TagStyle:        StyleCenter,
		PaddingChar:     "=",
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	// Tag that's too long to center (needs at least 4 chars for padding)
	entry.Message = "[VeryLongTag] Message"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	output := string(b)
	// Should fall back to non-centered when tag is too long
	if !strings.Contains(output, "[VeryLongTag]") {
		t.Errorf("Expected plain tag for long text, got: %s", output)
	}
}

func TestNewCustomFormatter_CenterBracketsDefaults(t *testing.T) {
	f := NewCustomFormatter()
	if f.TagStyle != StyleDefault {
		t.Error("Expected TagStyle to default to StyleDefault")
	}
	if f.PaddingChar != "•" {
		t.Errorf("Expected default PaddingChar '•', got '%s'", f.PaddingChar)
	}
}

func TestWithTagStyle_Option(t *testing.T) {
	tests := []struct {
		name         string
		style        TagStyle
		char         string
		expectedChar string
	}{
		{"StyleDefault", StyleDefault, "=", "="},
		{"StyleCenter", StyleCenter, "-", "-"},
		{"StyleRight", StyleRight, "·", "·"},
		{"Empty char defaults to bullet", StyleCenter, "", "•"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewCustomFormatter(WithTagStyle(tt.style, tt.char))
			if f.TagStyle != tt.style {
				t.Errorf("Expected TagStyle %v, got %v", tt.style, f.TagStyle)
			}
			if f.PaddingChar != tt.expectedChar {
				t.Errorf("Expected PaddingChar '%s', got '%s'", tt.expectedChar, f.PaddingChar)
			}
		})
	}
}

func TestFormatter_TagStyle_Default(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   false,
		TagStyle:        StyleDefault,
		PaddingChar:     "•",
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[Init] Starting"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	output := string(b)
	// Should have standard left-aligned bracket
	if !strings.Contains(output, "[Init]") {
		t.Errorf("Expected '[Init]', got: %s", output)
	}
	// Should NOT have padding dots
	if strings.Contains(output, "•") {
		t.Errorf("StyleDefault should not have padding dots, got: %s", output)
	}
}

func TestFormatter_TagStyle_Center(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   false,
		TagStyle:        StyleCenter,
		PaddingChar:     "•",
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[Init] Starting"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	output := string(b)
	// Should have centered padding with dots
	if !strings.Contains(output, "• Init •") {
		t.Errorf("Expected centered tag with '•' padding, got: %s", output)
	}
}

func TestFormatter_TagStyle_Right(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   false,
		TagStyle:        StyleRight,
		PaddingChar:     "•",
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[Init] Starting"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	output := string(b)
	// Should have right-aligned padding: [Init •...]
	if !strings.Contains(output, "[Init]•") {
		t.Errorf("Expected right-aligned tag with '•' padding, got: %s", output)
	}
}

func TestFormatter_TagStyle_Right_VariousLengths(t *testing.T) {
	tests := []struct {
		tag     string
		minDots int
	}{
		{"DB", 8},       // Very short: should have ~8 dots
		{"Init", 6},     // Short: should have ~6 dots
		{"Server", 4},   // Medium: should have ~4 dots
		{"Database", 2}, // Longer: should have ~2 dots
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			f := &CustomFormatter{
				UseColors:       false,
				CallerLevel:     logrus.WarnLevel,
				UseRelativePath: true,
				BracketPadding:  15,
				ColorBrackets:   false,
				TagStyle:        StyleRight,
				PaddingChar:     "•",
			}

			l := logrus.New()
			entry := l.WithField("test", "value")
			entry.Message = fmt.Sprintf("[%s] Message", tt.tag)
			entry.Level = logrus.InfoLevel

			b, err := f.Format(entry)
			if err != nil {
				t.Fatalf("Format error: %v", err)
			}

			output := string(b)
			// Verify right alignment
			if !strings.Contains(output, "["+tt.tag+"]•") {
				t.Errorf("Expected right-aligned tag '[%s •', got: %s", tt.tag, output)
			}
		})
	}
}

func TestFormatter_TagStyle_Right_WithColors(t *testing.T) {
	f := &CustomFormatter{
		UseColors:       true,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   true,
		TagStyle:        StyleRight,
		PaddingChar:     "·",
	}

	l := logrus.New()
	entry := l.WithField("test", "value")
	entry.Message = "[API] Request processed"
	entry.Level = logrus.InfoLevel

	b, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	output := string(b)
	// Should have yellow color
	if !strings.Contains(output, ColorYellow) {
		t.Error("Expected yellow color for brackets")
	}
	// Should have right-aligned padding with middle dots
	if !strings.Contains(stripANSI(output), "[API]·") {
		t.Errorf("Expected right-aligned tag with '·' padding, got: %s", output)
	}
	if !strings.Contains(output, ColorVeryDimGray) {
		t.Error("Expected dim gray padding color")
	}
}

func TestFormatter_TagStyle_AllStyles_Comparison(t *testing.T) {
	tag := "Auth"
	message := fmt.Sprintf("[%s] User logged in", tag)

	styles := []struct {
		name     string
		style    TagStyle
		contains string
	}{
		{"Default", StyleDefault, "[Auth]"},
		{"Center", StyleCenter, "• Auth •"},
		{"Right", StyleRight, "[Auth]•"},
	}

	for _, tt := range styles {
		t.Run(tt.name, func(t *testing.T) {
			f := &CustomFormatter{
				UseColors:       false,
				CallerLevel:     logrus.WarnLevel,
				UseRelativePath: true,
				BracketPadding:  15,
				ColorBrackets:   false,
				TagStyle:        tt.style,
				PaddingChar:     "•",
			}

			l := logrus.New()
			entry := l.WithField("test", "value")
			entry.Message = message
			entry.Level = logrus.InfoLevel

			b, err := f.Format(entry)
			if err != nil {
				t.Fatalf("Format error: %v", err)
			}

			output := string(b)
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.contains, output)
			}
		})
	}
}

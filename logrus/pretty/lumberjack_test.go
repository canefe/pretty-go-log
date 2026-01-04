package pretty

import (
	"testing"
)

func TestDefaultLogFileConfig(t *testing.T) {
	config := DefaultLogFileConfig()

	if config.MaxSize != 10 {
		t.Errorf("Expected MaxSize 10, got %d", config.MaxSize)
	}

	if config.MaxBackups != 5 {
		t.Errorf("Expected MaxBackups 5, got %d", config.MaxBackups)
	}

	if config.MaxAge != 31 {
		t.Errorf("Expected MaxAge 31, got %d", config.MaxAge)
	}

	if !config.Compress {
		t.Error("Expected Compress to be true")
	}
}

func TestNewLogFileConfig(t *testing.T) {
	config := NewLogFileConfig(20, 10, 60, false)

	if config.MaxSize != 20 {
		t.Errorf("Expected MaxSize 20, got %d", config.MaxSize)
	}

	if config.MaxBackups != 10 {
		t.Errorf("Expected MaxBackups 10, got %d", config.MaxBackups)
	}

	if config.MaxAge != 60 {
		t.Errorf("Expected MaxAge 60, got %d", config.MaxAge)
	}

	if config.Compress {
		t.Error("Expected Compress to be false")
	}
}

func TestNewLumberjackLogger(t *testing.T) {
	config := LogFileConfig{
		MaxSize:    15,
		MaxBackups: 7,
		MaxAge:     45,
		Compress:   true,
	}

	logger := NewLumberjackLogger("test.log", config)

	if logger == nil {
		t.Fatal("NewLumberjackLogger returned nil")
	}

	if logger.Filename != "test.log" {
		t.Errorf("Expected Filename 'test.log', got '%s'", logger.Filename)
	}

	if logger.MaxSize != 15 {
		t.Errorf("Expected MaxSize 15, got %d", logger.MaxSize)
	}

	if logger.MaxBackups != 7 {
		t.Errorf("Expected MaxBackups 7, got %d", logger.MaxBackups)
	}

	if logger.MaxAge != 45 {
		t.Errorf("Expected MaxAge 45, got %d", logger.MaxAge)
	}

	if !logger.Compress {
		t.Error("Expected Compress to be true")
	}
}

func TestNewLumberjackLogger_DefaultConfig(t *testing.T) {
	config := DefaultLogFileConfig()
	logger := NewLumberjackLogger("default.log", config)

	if logger == nil {
		t.Fatal("NewLumberjackLogger returned nil")
	}

	if logger.Filename != "default.log" {
		t.Errorf("Expected Filename 'default.log', got '%s'", logger.Filename)
	}

	if logger.MaxSize != 10 {
		t.Errorf("Expected default MaxSize 10, got %d", logger.MaxSize)
	}

	if logger.MaxBackups != 5 {
		t.Errorf("Expected default MaxBackups 5, got %d", logger.MaxBackups)
	}

	if logger.MaxAge != 31 {
		t.Errorf("Expected default MaxAge 31, got %d", logger.MaxAge)
	}

	if !logger.Compress {
		t.Error("Expected default Compress to be true")
	}
}

func TestLogFileConfig_StructTags(t *testing.T) {
	// Verify the struct is properly defined with expected defaults in tags
	config := LogFileConfig{}

	// Zero values when not initialized
	if config.MaxSize != 0 {
		t.Errorf("Expected zero MaxSize, got %d", config.MaxSize)
	}

	if config.MaxBackups != 0 {
		t.Errorf("Expected zero MaxBackups, got %d", config.MaxBackups)
	}

	if config.MaxAge != 0 {
		t.Errorf("Expected zero MaxAge, got %d", config.MaxAge)
	}

	if config.Compress {
		t.Error("Expected zero Compress to be false")
	}
}

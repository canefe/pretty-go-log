package pretty

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// Pre-compiled regexes for performance
var (
	bracketRegex     = regexp.MustCompile(`\[(.*?)\]`)
	parenthesesRegex = regexp.MustCompile(`\((.*?)\)`)
)

// TagStyle defines how bracketed tags are formatted
type TagStyle int

const (
	StyleDefault TagStyle = iota // [Auth]          - Standard left-aligned
	StyleCenter                  // [•• Auth ••]    - Centered with decorative padding
	StyleRight                   // [Auth ••••]     - Right-aligned with decorative padding
)

type CustomFormatter struct {
	UseColors     bool
	ShowCaller    bool
	ShowTimestamp bool
	// CallerLevel defines the minimum logrus level to start showing caller info.
	// e.g., if set to WarnLevel (3), it shows for Warn, Error, Fatal, and Panic.
	// Logrus levels: Panic(0) < Fatal(1) < Error(2) < Warn(3) < Info(4) < Debug(5) < Trace(6)
	CallerLevel logrus.Level
	// UseRelativePath uses relative file paths instead of absolute paths in caller info
	// When true, shows paths relative to working directory (e.g., "pkg/main.go:42")
	// When false, shows absolute paths (e.g., "/Users/name/project/pkg/main.go:42")
	UseRelativePath bool
	// BracketPadding sets the maximum padding for bracketed tags (e.g., "[Server]")
	// Tags longer than this will be truncated. Default: 15
	BracketPadding int
	// ColorBrackets enables or disables colored highlighting of bracketed tags
	// When false, brackets are shown without color even if UseColors is true
	ColorBrackets bool
	// TagStyle defines the bracket formatting style (Default, Center, or Right)
	// StyleDefault: [Auth]          - Standard left-aligned
	// StyleCenter:  [•• Auth ••]    - Centered with decorative padding
	// StyleRight:   [Auth ••••]     - Right-aligned with decorative padding
	TagStyle TagStyle
	// PaddingChar defines the character used for tag decoration
	// Common choices: "=", "-", "·", "•". Default: "•"
	PaddingChar string
}

// FormatterOption is a functional option for configuring CustomFormatter
type FormatterOption func(*CustomFormatter)

// WithColors enables or disables colored output
func WithColors(enabled bool) FormatterOption {
	return func(f *CustomFormatter) {
		f.UseColors = enabled
	}
}

// WithTimestamp enables or disables timestamp display
func WithTimestamp(enabled bool) FormatterOption {
	return func(f *CustomFormatter) {
		f.ShowTimestamp = enabled
	}
}

// WithCaller enables caller info and sets the minimum level to show it
func WithCaller(enabled bool, level logrus.Level) FormatterOption {
	return func(f *CustomFormatter) {
		f.ShowCaller = enabled
		f.CallerLevel = level
	}
}

// WithRelativePath configures whether to use relative or absolute paths in caller info
func WithRelativePath(enabled bool) FormatterOption {
	return func(f *CustomFormatter) {
		f.UseRelativePath = enabled
	}
}

// WithBracketPadding sets the maximum padding for bracketed tags
func WithBracketPadding(padding int) FormatterOption {
	return func(f *CustomFormatter) {
		if padding < 0 {
			padding = 0
		}
		f.BracketPadding = padding
	}
}

// WithColorBrackets enables or disables colored highlighting of bracketed tags
func WithColorBrackets(enabled bool) FormatterOption {
	return func(f *CustomFormatter) {
		f.ColorBrackets = enabled
	}
}

// WithTagStyle sets the bracket formatting style
//
// Examples:
//   - StyleDefault: [Auth]          (standard left-aligned)
//   - StyleCenter:  [•• Auth ••]    (centered with padding)
//   - StyleRight:   [Auth ••••]     (right-aligned with padding)
//
// The char parameter defines the padding character (e.g., "=", "-", "·", "•")
func WithTagStyle(style TagStyle, char string) FormatterOption {
	return func(f *CustomFormatter) {
		f.TagStyle = style
		if char != "" {
			f.PaddingChar = char
		} else {
			f.PaddingChar = "•"
		}
	}
}

// WithCenterBrackets is deprecated. Use WithTagStyle(StyleCenter, char) instead.
//
// Kept for backward compatibility.
func WithCenterBrackets(enabled bool, char string) FormatterOption {
	return func(f *CustomFormatter) {
		if enabled {
			f.TagStyle = StyleCenter
		} else {
			f.TagStyle = StyleDefault
		}
		if char != "" {
			f.PaddingChar = char
		} else {
			f.PaddingChar = "•"
		}
	}
}

// NewCustomFormatter creates a new formatter with the given options
func NewCustomFormatter(opts ...FormatterOption) *CustomFormatter {
	// Defaults: colors on, timestamps off, caller on for Warn and above, relative paths, 15 char bracket padding, colored brackets, default tag style
	f := &CustomFormatter{
		UseColors:       true,
		ShowCaller:      true,
		ShowTimestamp:   false,
		CallerLevel:     logrus.WarnLevel,
		UseRelativePath: true,
		BracketPadding:  15,
		ColorBrackets:   true,
		TagStyle:        StyleDefault,
		PaddingChar:     "•",
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

const (
	ColorReset       = "\033[0m"
	ColorRed         = "\033[31m"
	ColorGreen       = "\033[32m"
	ColorYellow      = "\033[33m"
	ColorMagenta     = "\033[35m"
	ColorCyan        = "\033[36m"
	ColorGray        = "\033[90m"
	ColorDimGray     = "\033[2;37m"
	ColorVeryDimGray = "\033[38;5;242m"
	ColorDarkGray    = "\033[38;5;240m"
	ColorDeepGray    = "\033[38;5;238m"
	ColorAbyssalGray = "\033[38;5;236m"
	ColorGhostGray   = "\033[38;5;234m"
)

// formatCommon formats common log elements (timestamp, level, message, colors)
func (f *CustomFormatter) formatCommon(entry *logrus.Entry) (timestamp, level, message, colorCode, resetCode string) {
	timestamp = entry.Time.Format("2006-01-02 15:04:05")
	level = strings.ToUpper(entry.Level.String())

	if level == "WARNING" {
		level = "WARN"
	}

	message = entry.Message

	if f.UseColors {
		switch entry.Level {
		case logrus.PanicLevel:
			colorCode = ColorMagenta
		case logrus.TraceLevel:
			colorCode = ColorCyan
		case logrus.DebugLevel:
			colorCode = ColorCyan
		case logrus.InfoLevel:
			colorCode = ColorGreen
		case logrus.WarnLevel:
			colorCode = ColorYellow
		case logrus.ErrorLevel:
			colorCode = ColorRed
		case logrus.FatalLevel:
			colorCode = ColorMagenta
		default:
			colorCode = ColorReset
		}
		resetCode = ColorReset
	}

	return timestamp, level, message, colorCode, resetCode
}

// formatCallerInfo formats the caller information with optional coloring
func (f *CustomFormatter) formatCallerInfo(entry *logrus.Entry) string {
	if entry.Caller == nil {
		return ""
	}

	filePath := entry.Caller.File

	// Use absolute or relative path based on configuration
	if !f.UseRelativePath {
		// Get absolute path
		absFile, err := filepath.Abs(filePath)
		if err == nil {
			filePath = absFile
		}
	} else {
		// Get relative path from current working directory
		if wd, err := os.Getwd(); err == nil {
			if relPath, err := filepath.Rel(wd, filePath); err == nil {
				filePath = relPath
			}
		}
	}

	callerInfo := fmt.Sprintf("└─ at (%s:%d)", filePath, entry.Caller.Line)

	// Apply color to parentheses content if colors are enabled
	if f.UseColors {
		callerInfo = parenthesesRegex.ReplaceAllStringFunc(callerInfo, func(bracketed string) string {
			return fmt.Sprintf("%s%s%s", ColorVeryDimGray, bracketed, ColorReset)
		})
	}

	return callerInfo
}

func (f *CustomFormatter) coloredTagWithPadding(inner string, maxPadding int, style TagStyle) string {
	tagColor := ColorYellow
	padColor := ColorVeryDimGray

	fill := f.PaddingChar
	if fill == "" {
		fill = "•"
	}

	switch style {
	case StyleCenter:
		availableSpace := maxPadding - 2
		if len(inner) >= availableSpace-2 {
			return tagColor + "[" + inner + "]" + ColorReset
		}
		totalDots := availableSpace - len(inner) - 2
		leftDots := totalDots / 2
		rightDots := totalDots - leftDots

		var b strings.Builder
		b.WriteString(tagColor)
		b.WriteByte('[')
		b.WriteString(padColor)
		b.WriteString(strings.Repeat(fill, leftDots))
		b.WriteString(tagColor)
		b.WriteByte(' ')
		b.WriteString(inner)
		b.WriteByte(' ')
		b.WriteString(padColor)
		b.WriteString(strings.Repeat(fill, rightDots))
		b.WriteString(tagColor)
		b.WriteByte(']')
		b.WriteString(ColorReset)
		return b.String()
	case StyleRight:
		availableSpace := maxPadding - 2
		if len(inner) >= availableSpace-1 {
			return tagColor + "[" + inner + "]" + ColorReset
		}
		totalDots := availableSpace - len(inner) - 1

		var b strings.Builder
		b.WriteString(tagColor)
		b.WriteByte('[')
		b.WriteString(inner)
		b.WriteByte(']')
		b.WriteString(padColor)
		b.WriteString(strings.Repeat(fill, totalDots))
		b.WriteString(tagColor)
		b.WriteByte(' ')
		b.WriteString(ColorReset)
		return b.String()
	default:
		return tagColor + "[" + inner + "]" + ColorReset
	}
}

// centerTag centers a tag within the available space with decorative padding
//
// Example: "Auth" -> "[•• Auth ••]" (assuming maxPadding=15)
func (f *CustomFormatter) centerTag(inner string, maxPadding int) string {
	availableSpace := maxPadding - 2 // Subtract 2 for the brackets
	if len(inner) >= availableSpace-2 {
		return "[" + inner + "]" // Not enough space, return as-is
	}

	totalDots := availableSpace - len(inner) - 2 // 2 spaces around the text
	leftDots := totalDots / 2
	rightDots := totalDots - leftDots

	fill := f.PaddingChar
	if fill == "" {
		fill = "•"
	}

	return fmt.Sprintf("[%s %s %s]",
		strings.Repeat(fill, leftDots),
		inner,
		strings.Repeat(fill, rightDots))
}

// rightPadTag right-aligns a tag within the available space with decorative padding
//
// Example: "Auth" -> "[Auth]••••" (assuming maxPadding=15)
func (f *CustomFormatter) rightPadTag(inner string, maxPadding int) string {
	availableSpace := maxPadding - 2 // Subtract 2 for the brackets
	if len(inner) >= availableSpace-1 {
		return "[" + inner + "]" // Not enough space, return as-is
	}

	totalDots := availableSpace - len(inner) - 1 // 1 space before the dots
	fill := f.PaddingChar
	if fill == "" {
		fill = "•"
	}

	return fmt.Sprintf("[%s]%s ", inner, strings.Repeat(fill, totalDots))
}

// appendFields sorts and appends structured data to the log line
func (f *CustomFormatter) appendFields(b *strings.Builder, data logrus.Fields) {
	// 1. Sort the keys for consistent output across runs
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. Append each field
	b.WriteByte(' ') // Lead with a space to separate from the message
	for i, k := range keys {
		v := data[k]

		if f.UseColors {
			// Key in Dim Gray, Value in a slightly brighter Gray
			b.WriteString(fmt.Sprintf("%s%s=%s%v%s", ColorVeryDimGray, k, ColorGray, v, ColorReset))
		} else {
			b.WriteString(fmt.Sprintf("%s=%v", k, v))
		}

		// Add space between fields, but not after the last one
		if i < len(keys)-1 {
			b.WriteByte(' ')
		}
	}
}

func stripANSI(str string) string {
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansi.ReplaceAllString(str, "")
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp, level, message, colorCode, resetCode := f.formatCommon(entry)
	var b strings.Builder

	// 1. Timestamp & Level
	if f.ShowTimestamp {
		b.WriteString("[" + timestamp + "] ")
	}
	b.WriteString(colorCode + level)
	b.WriteString(strings.Repeat(" ", max(0, 6-len(level))) + resetCode + " ")

	// 2. Bracketed Tag Handling
	maxPadding := f.BracketPadding
	if maxPadding <= 0 {
		maxPadding = 15
	}

	var displayTag string
	if loc := bracketRegex.FindStringIndex(message); loc != nil {
		inner := strings.Trim(message[loc[0]:loc[1]], "[]")

		// 1. Generate the styled tag string
		if f.UseColors && f.ColorBrackets {
			if f.TagStyle == StyleCenter || f.TagStyle == StyleRight {
				displayTag = f.coloredTagWithPadding(inner, maxPadding, f.TagStyle)
			} else {
				displayTag = ColorYellow + "[" + inner + "]" + ColorReset
			}
		} else {
			switch f.TagStyle {
			case StyleCenter:
				displayTag = f.centerTag(inner, maxPadding)
			case StyleRight:
				displayTag = f.rightPadTag(inner, maxPadding)
			default:
				displayTag = "[" + inner + "]"
			}
		}

		// 2. Write the tag and clean up the message
		b.WriteString(displayTag)
		message = strings.TrimSpace(message[:loc[0]] + message[loc[1]:])
	}

	// 3. Calculate remaining gutter space
	// stripANSI ensures we don't count invisible color codes
	visibleLen := len(stripANSI(displayTag))
	if visibleLen < maxPadding {
		b.WriteString(strings.Repeat(" ", maxPadding-visibleLen))
	}
	b.WriteByte(' ') // Single space separator before the message text

	// 3. Message & Fields
	b.WriteString(message)
	if len(entry.Data) > 0 {
		f.appendFields(&b, entry.Data)
	}

	// 4. Caller Info
	if f.ShowCaller && entry.Level <= f.CallerLevel {
		// Calculate how many spaces we need to skip to reach the message column
		// Timestamp (approx 22) + Level (7) + Gutter (maxPadding + 1)
		prefixWidth := 0
		if f.ShowTimestamp {
			prefixWidth += 22 // "[2006-01-02 15:04:05] "
		}
		prefixWidth += 7          // "LEVEL  " (Level 6 + 1 space)
		prefixWidth += maxPadding // The tag gutter
		prefixWidth += 1          // The final separator space

		indent := strings.Repeat(" ", prefixWidth)
		b.WriteString("\n" + indent + f.formatCallerInfo(entry))
	}

	b.WriteByte('\n')
	return []byte(b.String()), nil
}

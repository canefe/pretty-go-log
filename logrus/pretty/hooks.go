package pretty

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type MultiWriterWithFormattersConfig struct {
	format       FormatType
	showCaller   bool
	customFormat *CustomFormatter
}
type writerPair struct {
	w io.Writer
	f logrus.Formatter
}

type MultiWriter struct {
	pairs []writerPair
	cfg   MultiWriterWithFormattersConfig
}

func NewMultiWriter(cfg MultiWriterWithFormattersConfig) *MultiWriter {
	return &MultiWriter{cfg: cfg}
}

func (mw *MultiWriter) AddWriter(w io.Writer, useColors, showTime bool) {
	var f logrus.Formatter

	if mw.cfg.customFormat != nil {
		custom := *mw.cfg.customFormat
		custom.UseColors = useColors
		custom.ShowTimestamp = showTime
		f = &custom
	} else {
		switch mw.cfg.format {
		case FormatJSON:
			f = &logrus.JSONFormatter{}
		case FormatPlain:
			f = &CustomFormatter{
				UseColors:       useColors,
				ShowCaller:      mw.cfg.showCaller,
				ShowTimestamp:   showTime,
				CallerLevel:     logrus.WarnLevel,
				UseRelativePath: true,
				BracketPadding:  15,
				ColorBrackets:   true,
			}
		default:
			f = &logrus.TextFormatter{ForceColors: useColors}
		}
	}

	mw.pairs = append(mw.pairs, writerPair{w, f})
}

func (mw *MultiWriter) WriteEntry(e *logrus.Entry) error {
	for _, p := range mw.pairs {
		buf, err := p.f.Format(e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "log format err: %v\n", err)
			continue
		}
		if _, err := p.w.Write(buf); err != nil {
			fmt.Fprintf(os.Stderr, "log write err: %v\n", err)
		}
	}
	return nil
}

type CustomHook struct {
	mw *MultiWriter
}

func (h *CustomHook) Levels() []logrus.Level     { return logrus.AllLevels }
func (h *CustomHook) Fire(e *logrus.Entry) error { return h.mw.WriteEntry(e) }

package logging

import (
	"bytes"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type ColorScheme struct {
	InfoLevel  string
	WarnLevel  string
	ErrorLevel string
	FatalLevel string
	PanicLevel string
	DebugLevel string
	Timestamp  string
}

type compiledColorScheme struct {
	InfoLevel  func(string) string
	WarnLevel  func(string) string
	ErrorLevel func(string) string
	FatalLevel func(string) string
	PanicLevel func(string) string
	DebugLevel func(string) string
	Timestamp  func(string) string
}

const defaultTimestampFormat = time.RFC3339

var (
	defaultLogForamt = "[%level%]: %time% - [%trace_id%] %msg%"
	//baseTimestamp time.Time = time.Now()
	defaultColorScheme = &ColorScheme{
		InfoLevel:  "green",
		WarnLevel:  "yellow",
		ErrorLevel: "red",
		FatalLevel: "red",
		PanicLevel: "red",
		DebugLevel: "blue",
		Timestamp:  "black+h",
	}
	defaultCompiledColorScheme = compileColorScheme(defaultColorScheme)
)

type TextFormatter struct {
	LogFormat       string
	TimestampFormat string
	colorScheme     *compiledColorScheme
	QuoteCharacter  string
	isTerminal      bool
	sync.Once
}

func getCompiledColor(main string, fallback string) func(string) string {
	var style string
	if main != "" {
		style = main
	} else {
		style = fallback
	}
	return ansi.ColorFunc(style)
}

func compileColorScheme(s *ColorScheme) *compiledColorScheme {
	return &compiledColorScheme{
		InfoLevel:  getCompiledColor(s.InfoLevel, defaultColorScheme.InfoLevel),
		WarnLevel:  getCompiledColor(s.WarnLevel, defaultColorScheme.WarnLevel),
		ErrorLevel: getCompiledColor(s.ErrorLevel, defaultColorScheme.ErrorLevel),
		FatalLevel: getCompiledColor(s.FatalLevel, defaultColorScheme.FatalLevel),
		PanicLevel: getCompiledColor(s.PanicLevel, defaultColorScheme.PanicLevel),
		DebugLevel: getCompiledColor(s.DebugLevel, defaultColorScheme.DebugLevel),
		Timestamp:  getCompiledColor(s.Timestamp, defaultColorScheme.Timestamp),
	}
}

// 初始化参数
func (f *TextFormatter) init(entry *logrus.Entry) {
	if len(f.QuoteCharacter) == 0 {
		f.QuoteCharacter = "\""
	}

	if entry.Logger != nil {
		f.isTerminal = f.checkIfTerminal(entry.Logger.Out)
	}
}

// 检测是否是终端
func (f *TextFormatter) checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func (f *TextFormatter) SetColorScheme(colorScheme *ColorScheme) {
	f.colorScheme = compileColorScheme(colorScheme)
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	f.Do(func() {
		f.init(entry)
	})

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	timestampFormat = "2006-01-02 15:04:05"
	output := f.LogFormat
	if output == "" {
		output = defaultLogForamt
	}

	timestamp := entry.Time.Format(timestampFormat)

	var willReplacer []string
	if f.isTerminal {
		willReplacer = f.outputTerminal(entry, timestamp)
	} else {
		willReplacer = f.output(entry, timestamp)
	}

	r := strings.NewReplacer(willReplacer...)

	output = r.Replace(output)
	var b bytes.Buffer
	b.WriteString(output)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) output(entry *logrus.Entry, timestamp string) []string {
	level := entry.Level.String()
	msg := entry.Message
	traceId := new(TraceIdGenerate).get()
	willReplacer := []string{"%level%", level, "%trace_id%", traceId, "%msg%", msg, "%time%", timestamp}
	for key, value := range entry.Data {
		if s, ok := value.(string); ok {
			willReplacer = append(willReplacer, "%"+key+"%", s)
		}
	}
	return willReplacer
}

func (f *TextFormatter) outputTerminal(entry *logrus.Entry, timestamp string) []string {
	var levelColor func(string) string
	var colorScheme *compiledColorScheme
	if f.colorScheme == nil {
		colorScheme = defaultCompiledColorScheme
	} else {
		colorScheme = f.colorScheme
	}

	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = colorScheme.InfoLevel
	case logrus.WarnLevel:
		levelColor = colorScheme.WarnLevel
	case logrus.ErrorLevel:
		levelColor = colorScheme.ErrorLevel
	case logrus.FatalLevel:
		levelColor = colorScheme.FatalLevel
	case logrus.PanicLevel:
		levelColor = colorScheme.PanicLevel
	default:
		levelColor = colorScheme.DebugLevel
	}
	level := entry.Level.String()
	msg := entry.Message
	traceId := new(TraceIdGenerate).get()
	willReplacer := []string{"%level%", levelColor(level), "%trace_id%", traceId, "%msg%", msg, "%time%", colorScheme.Timestamp(timestamp)}
	for key, value := range entry.Data {
		if s, ok := value.(string); ok {
			willReplacer = append(willReplacer, "%"+key+"%", s)
		}
	}
	return willReplacer
}

package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type RootFields struct {
	Timestamp string
	Func      string
	Level     logrus.Level
	Fields    interface{}
}

type Formatter struct {
	PrettyPrint   bool
	CustomCaption string
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	levelColor := getColorByLevel(entry.Level)
	root := RootFields{Timestamp: entry.Time.Format("2006-01-02 15:04:05"), Level: entry.Level,
		//CustomCaption: entry.CustomCaption, // not possible in logrus...
		Fields: encode(entry.Message)}

	b.WriteString(root.Timestamp)

	b.WriteString(" - ")
	_, _ = fmt.Fprintf(b, "\x1b[%d;1m", levelColor)
	b.WriteString(strings.ToUpper(root.Level.String()))
	b.WriteString("\x1b[0m")
	b.WriteString(getPaddingByLevel(entry.Level))
	b.WriteString(" - ")

	// if entry.HasCaller() {
	// 	caller := getCaller(entry.Caller)
	// 	fc := caller.Function
	// 	file := fmt.Sprintf("%s:%d", caller.File, caller.Line)
	// 	b.WriteString(prettierCaller(file, fc))
	// }

	b.WriteString(fmt.Sprintf("[%s]", f.CustomCaption))
	b.WriteString(fmt.Sprintf("%16s", " - "))

	// _, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)

	var data string
	data = marshal(root.Fields)

	b.WriteString(data)
	b.WriteString("\x1b[0m")

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func prettierCaller(file string, function string) string {
	dirs := strings.Split(file, "/")
	fileDesc := strings.Join(dirs[len(dirs)-2:], "/")

	funcs := strings.Split(function, ".")
	funcDesc := strings.Join(funcs[len(funcs)-2:], ".")

	return "[" + fileDesc + ":" + funcDesc + "]"
}

func encode(message string) interface{} {
	if data := encodeForJsonString(message); data != nil {
		return data
	} else {
		return message
	}
}

func encodeForJsonString(message string) map[string]interface{} {
	// jsonstring
	inInterface := make(map[string]interface{})
	if err := json.Unmarshal([]byte(message), &inInterface); err != nil {
		//fmt.Print("err !!!! " , err.Error())
		return nil
	}
	return inInterface
}

const (
	colorRed      = 31
	colorGreen    = 32
	colorYellow   = 33
	colorDarkBlue = 34
	colorBlue     = 36
	colorGray     = 37
)

func getPaddingByLevel(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return fmt.Sprintf("%3s", "")
	case logrus.DebugLevel:
		return fmt.Sprintf("%3s", "")
	case logrus.InfoLevel:
		return fmt.Sprintf("%4s", "")
	case logrus.WarnLevel:
		return fmt.Sprintf("%1s", "")
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return fmt.Sprintf("%3s", "")
	default:
		return fmt.Sprintf("%2s", "")
	}
}

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.TraceLevel:
		return colorGray
	case logrus.DebugLevel:
		return colorBlue
	case logrus.InfoLevel:
		return colorGreen
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorDarkBlue
	}
}

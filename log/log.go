package log

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"linkshortener/lib/tool"
	"linkshortener/setting"
	"os"
	"strings"
	"time"
)

var (
	Stdout    *os.File
	NullOut   *os.File
	logWriter io.Writer
)

func IsDebug() bool {
	return setting.Cfg.LOG.Debug
}

func InitLog() {
	var err error
	timeStr := tool.Now()
	Stdout = os.Stdout
	NullOut, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0600)

	color.Set(color.FgMagenta)
	defer color.Unset()

	path := "./logs/"
	_ = tool.Mkdir(path)
	logWriter, err = os.OpenFile(path+tool.NowDay()+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		_, _ = fmt.Fprintf(Stdout, "[PANIC] ["+timeStr+"] Error opening log file: %s\n", err)
		_, _ = fmt.Fprintf(Stdout, "[PANIC] ["+timeStr+"] Program Exit after 5 Second\n")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	if setting.Cfg.RunMode != "dev" {
		os.Stdout = NullOut
		os.Stderr = NullOut
	}
}

func Close() {
	if logWriter != nil {
		if closer, ok := logWriter.(io.Closer); ok {
			_ = closer.Close()
		}
		logWriter = nil
	}
}

func DebugPrint(format string, values ...interface{}) {
	if IsDebug() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		timeStr := tool.Now()
		_, _ = fmt.Fprintf(logWriter, "[DEBUG] ["+timeStr+"] "+format, values...)
		_, _ = fmt.Fprintf(Stdout, "[DEBUG] ["+timeStr+"] "+format, values...)
	}
}

func InfoPrint(format string, values ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	timeStr := tool.Now()
	_, _ = fmt.Fprintf(logWriter, "[INFO]  ["+timeStr+"] "+format, values...)
	_, _ = fmt.Fprintf(Stdout, "[INFO]  ["+timeStr+"] "+format, values...)
}

func WarnPrint(format string, values ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	timeStr := tool.Now()

	color.Set(color.FgYellow)
	defer color.Unset()

	_, _ = fmt.Fprintf(logWriter, "[WARN]  ["+timeStr+"] "+format, values...)
	_, _ = fmt.Fprintf(Stdout, "[WARN]  ["+timeStr+"] "+format, values...)
}

func ErrorPrint(format string, values ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	timeStr := tool.Now()

	color.Set(color.FgRed)
	defer color.Unset()

	_, _ = fmt.Fprintf(logWriter, "[ERROR] ["+timeStr+"] "+format, values...)
	_, _ = fmt.Fprintf(Stdout, "[ERROR] ["+timeStr+"] "+format, values...)
}

func PanicPrint(format string, values ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	timeStr := tool.Now()

	color.Set(color.FgMagenta)
	defer color.Unset()

	_, _ = fmt.Fprintf(logWriter, "[PANIC] ["+timeStr+"] "+format, values...)
	_, _ = fmt.Fprintf(Stdout, "[PANIC] ["+timeStr+"] "+format, values...)
	WarnPrint("Program Exit after 5 Second")
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func Errorf(format string, values ...interface{}) error {
	ErrorPrint(format, values...)
	return fmt.Errorf(format, values...)
}

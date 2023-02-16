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

func IsDebug() bool {
	return setting.Cfg.LOG.Debug
}

func GetWriter() io.Writer {
	var w io.Writer
	path := "./logs/"
	_ = tool.Mkdir(path)
	setting.Cfg.LOG.File, _ = os.OpenFile(path+tool.NowDay()+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	w = setting.Cfg.LOG.File
	return w
}

func Close() {
	_ = setting.Cfg.LOG.File.Close()
}

func DebugPrint(format string, values ...interface{}) {
	if IsDebug() {
		w := GetWriter()
		defer Close()
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		timeStr := tool.Now()
		_, _ = fmt.Fprintf(w, "[DEBUG] ["+timeStr+"] "+format, values...)
		_, _ = fmt.Fprintf(os.Stdout, "[DEBUG] ["+timeStr+"] "+format, values...)
	}
}

func InfoPrint(format string, values ...interface{}) {
	w := GetWriter()
	defer Close()
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	timeStr := tool.Now()
	_, _ = fmt.Fprintf(w, "[INFO]  ["+timeStr+"] "+format, values...)
	_, _ = fmt.Fprintf(os.Stdout, "[INFO]  ["+timeStr+"] "+format, values...)
}

func WarnPrint(format string, values ...interface{}) {
	w := GetWriter()
	defer Close()
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	timeStr := tool.Now()

	color.Set(color.FgYellow)
	defer color.Unset()

	_, _ = fmt.Fprintf(w, "[WARN]  ["+timeStr+"] "+format, values...)
	_, _ = fmt.Fprintf(os.Stdout, "[WARN]  ["+timeStr+"] "+format, values...)
}

func ErrorPrint(format string, values ...interface{}) {
	w := GetWriter()
	defer Close()
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	timeStr := tool.Now()

	color.Set(color.FgRed)
	defer color.Unset()

	_, _ = fmt.Fprintf(w, "[ERROR] ["+timeStr+"] "+format, values...)
	_, _ = fmt.Fprintf(os.Stdout, "[ERROR] ["+timeStr+"] "+format, values...)
}

func PanicPrint(format string, values ...interface{}) {
	w := GetWriter()
	defer Close()
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	timeStr := tool.Now()

	color.Set(color.FgMagenta)
	defer color.Unset()

	_, _ = fmt.Fprintf(w, "[PANIC] ["+timeStr+"] "+format, values...)
	_, _ = fmt.Fprintf(os.Stdout, "[PANIC] ["+timeStr+"] "+format, values...)
	WarnPrint("Program Exit after 5 Second")
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

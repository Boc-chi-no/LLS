package setting

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ini/ini"
	"github.com/rakyll/statik/fs"
	"linkshortener/lib/tool"
	"linkshortener/model"
	"os"
)

var Cfg model.Config

func InitSetting() {
	var err error

	if !tool.FileExist("./app.ini") {
		sfs, _ := fs.New()
		settingBytes, _ := fs.ReadFile(sfs, "/statik/app.ini")
		w, _ := os.OpenFile("./app.ini", os.O_WRONLY|os.O_CREATE, 0666)
		_, _ = w.Write(settingBytes)
		_ = w.Close()
		color.Set(color.FgYellow)
		defer color.Unset()
		_, _ = fmt.Fprintf(os.Stdout, "[WARN]  ["+tool.Now()+"] The configuration file does not exist and has been automatically generated\n")
	}

	cfgFile, err := ini.Load("app.ini")
	if err != nil {
		color.Set(color.FgMagenta)
		defer color.Unset()
		_, _ = fmt.Fprintf(os.Stdout, "[PANIC]  ["+tool.Now()+"] Fail to Load ‘app.ini’: %s", err, "\n")
		os.Exit(0)
	}
	err = cfgFile.MapTo(&Cfg)
	if err != nil {
		color.Set(color.FgMagenta)
		defer color.Unset()
		_, _ = fmt.Fprintf(os.Stdout, "[PANIC]  ["+tool.Now()+"] Fail to Map ‘app.ini’: %s", err, "\n")
		os.Exit(0)
	}
}

//go:generate statik -f -p=fs -src=./static/
//go:generate go fmt ./fs/statik.go

package main

import (
	"linkshortener/controller"
	"linkshortener/db"
	"linkshortener/fs"
	"linkshortener/log"
	"linkshortener/setting"
	"time"
)

func main() {
	setting.InitSetting()
	log.InitLog()
	fs.InitFs()
	fs.InitFont()
	fs.InitUap()
	fs.InitIPData()
	fs.InitI18n()

	db.InitDB()
	db.InitModel()

	controller.InitController()
	controller.InitRouter()

	controller.RunServer()

	time.Sleep(5 * time.Second)
}

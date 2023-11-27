//go:generate statik -f -src=./resources
//go:generate go fmt statik/statik.go

package main

import (
	"linkshortener/controller"
	"linkshortener/db"
	"linkshortener/log"
	"linkshortener/setting"
	_ "linkshortener/statik"
	"linkshortener/statikFS"
	"time"
)

func main() {
	setting.InitSetting()
	log.InitLog()
	statikFS.InitFs()
	statikFS.InitFont()
	statikFS.InitUap()
	statikFS.InitIPData()
	statikFS.InitI18n()

	db.InitDB()
	db.InitModel()

	controller.InitController()
	controller.InitRouter()

	controller.RunServer()

	time.Sleep(5 * time.Second)
}

//go:generate statik -f -src=./public
//go:generate go fmt statik/statik.go

package main

import (
	"linkshortener/controller"
	"linkshortener/db"
	"linkshortener/setting"
	_ "linkshortener/statik"
	"linkshortener/statikFS"
	"time"
)

func main() {
	setting.InitSetting()
	statikFS.InitFs()
	statikFS.InitFont()
	statikFS.InitUap()
	statikFS.InitIPData()

	db.InitDB()
	db.InitModel()

	controller.InitController()
	controller.InitRouter()

	controller.RunServer()

	time.Sleep(5 * time.Second)
}

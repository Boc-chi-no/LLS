package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"linkshortener/db"
	"linkshortener/lib/ip2location"
	"linkshortener/lib/uap"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
	"time"
)

// Redirect This method performs the redirection of the shortened link.
// Usage:
// Just hit http://localhost:8040/s/4nGHqG ( use the generated hash )
func Redirect(c *gin.Context) {
	shortHash := c.Param("hash")

	var res []model.Link
	table:=db.NewModel(setting.Cfg.MongoDB.Database, "links")


	table.Find(bson.D{{Key:"_id",Value: shortHash},{Key:"delete",Value: false}},&res)

	if res != nil&&len(res) > 0 {
		go accessLogWorker(c.ClientIP(),shortHash,c.Request.Header,time.Now().Unix())
		log.DebugPrint("RedirectLink: %s",res[0].URL)
		c.Redirect(http.StatusTemporaryRedirect, res[0].URL)
	} else {
		model.FailureResponse(c,404,404,"未找到查询的链接!","")
	}
}

func accessLogWorker(ip string,hash string,header http.Header,nowTime int64) {
	qqWry := ip2location.NewQQwry()
	location := qqWry.Find(ip)
	uaInfo := uap.Parse(header)

	var linkInfo = model.LinkInfo{
		Hash: hash,
		IP: ip,
		Header: header,
		Location:location,
		UAInfo: uaInfo,
		Created: nowTime,
	}

	table:=db.NewModel(setting.Cfg.MongoDB.Database, "link_access")
	res:=table.InsertOne(linkInfo)
	if res==nil{
		log.WarnPrint("Failed to write access log to database!")
	}
}
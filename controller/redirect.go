package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"linkshortener/db"
	"linkshortener/i18n"
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
	req := model.RedirectLinkReq{}
	localizer := i18n.GetLocalizer(c)

	if err := c.BindUri(&req); err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("deserializationFailed", nil), "")
		log.ErrorPrint("Deserialization failed: %s", err)
		return
	}

	var res []model.Link
	table := db.SetModel(setting.Cfg.MongoDB.Database, "links")
	_ = table.Find(bson.D{{Key: "_id", Value: req.Hash}, {Key: "delete", Value: false}}, &res)

	if res != nil && len(res) > 0 {
		go accessLogWorker(c.ClientIP(), req.Hash, c.Request.Header, time.Now().Unix())
		log.DebugPrint("RedirectLink: %s", res[0].URL)
		c.Redirect(http.StatusTemporaryRedirect, res[0].URL)
	} else {
		model.FailureResponse(c, 404, 404, localizer.GetMessage("noLinkFound", nil), "")
	}
}

func accessLogWorker(ip string, hash string, header http.Header, nowTime int64) {
	qqWry := ip2location.NewQQwry()
	location := qqWry.Find(ip)
	uaInfo := uap.Parse(header)

	var linkInfo = model.LinkInfo{
		Hash:     hash,
		IP:       ip,
		Header:   header,
		Location: location,
		UAInfo:   uaInfo,
		Created:  nowTime,
	}

	table := db.SetModel(setting.Cfg.MongoDB.Database, "link_access")
	res, _ := table.InsertOne(linkInfo)
	if res == nil {
		log.WarnPrint("Failed to write access log to database!")
	}
}

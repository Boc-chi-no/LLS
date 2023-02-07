package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"linkshortener/db"
	"linkshortener/lib/shorten"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
)

// InsertLink This method saves the redirection into the mongo database.
//
// Usage:
// To shorten an URL, just http POST to http://localhost:8040/api/generate_link with the following json payload (example):
//
// {
//    "src_url": "http://localhost:8040/",
//    "captcha": "5SDF1"
// }
func InsertLink(c *gin.Context) {
	var req model.InsertLinkReq

	if err := c.BindJSON(&req); err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, "序列化失败", "")
		log.ErrorPrint("Deserialization failed: %s", err)
		return
	}

	// 初始化session对象
	session := sessions.Default(c)
	sessionCaptcha := session.Get("captcha")
	session.Delete("captcha")
	_ = session.Save()

	if sessionCaptcha != req.CAPTCHA {
		model.FailureResponse(c, http.StatusForbidden, http.StatusForbidden, "验证码检验失败!", "")
		return
	}

	link := shorten.ShortenLink(req)

	table := db.NewModel(setting.Cfg.MongoDB.Database, "links")
	res := table.InsertOne(link)
	if res == nil {
		model.FailureResponse(c, 500, 500, "写入数据库失败!", "")
		return
	}

	log.DebugPrint("SrcLink: %s, ShortenLink: http://%s/s/%s", req.URL, setting.Cfg.HTTP.Listen, res.InsertedID)

	data := map[string]interface{}{
		"hash":  link.ShortHash,
		"token": link.Token,
	}

	model.SuccessResponse(c, data)
}

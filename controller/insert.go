package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"linkshortener/db"
	"linkshortener/i18n"
	"linkshortener/lib/shorten"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
	"strings"
)

// GenerateLink This method saves the redirection into the mongo database.
//
// Usage:
// To shorten a URL, just http POST to http://localhost:8040/api/generate_link with the following json payload (example):
//
//	{
//	   "src_url": "http://localhost:8040/",
//	   "captcha": "5SDF1"
//	}
func GenerateLink(c *gin.Context) {
	var req model.InsertLinkReq
	localizer := i18n.GetLocalizer(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("deserializationFailed", nil), "")
		log.ErrorPrint("Deserialization failed: %s", err)
		return
	}

	// Initialize session object
	session := sessions.Default(c)
	sessionCaptcha := session.Get("captcha")
	session.Delete("captcha")
	_ = session.Save()

	if sessionCaptcha != req.CAPTCHA {
		model.FailureResponse(c, http.StatusForbidden, http.StatusForbidden, localizer.GetMessage("captchaVerificationFailed", nil), "")
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("invalidUrl", nil), "")
		log.WarnPrint("Illegal URL: %s", req.URL)
		return
	}

	link := shorten.GenerateShortenLink(req)

	table := db.SetModel(setting.Cfg.MongoDB.Database, "links")
	res, _ := table.InsertOne(link)
	if res == nil {
		model.FailureResponse(c, http.StatusInternalServerError, http.StatusInternalServerError, localizer.GetMessage("databaseOperationFailed", nil), "")
		return
	}

	log.DebugPrint("SrcLink: %s, GenerateShortenLink: /s/%s", req.URL, res.InsertedID)

	data := map[string]interface{}{
		"hash":  link.ShortHash,
		"token": link.Token,
	}

	model.SuccessResponse(c, data)
}

package controller

import (
	"linkshortener/db"
	"linkshortener/i18n"
	"linkshortener/lib/shorten"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GenerateLink This method saves the redirection into the mongo database.
//
// Usage:
// To shorten a URL, just http POST to {BasePath}/api/generate_link with the following json payload (example):
//
//	{
//	   "src_url": "http://localhost:8040/",
//	   "captcha": "5SDF1"
//	}
func GenerateLink(c *gin.Context) {
	var req model.InsertLinkReq
	localizer := i18n.GetLocalizer(c)
	now := time.Now().Unix()

	if err := c.ShouldBindJSON(&req); err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("deserializationFailed", nil), err.Error())
		log.ErrorPrint("Deserialization failed: %s", err)
		return
	}

	// Initialize session object
	session := sessions.Default(c)
	sessionCaptcha := tool.SafeSessionGet(session, "captcha")
	session.Delete("captcha")
	_ = session.Save()

	if sessionCaptcha != req.CAPTCHA {
		model.FailureResponse(c, http.StatusForbidden, http.StatusForbidden, localizer.GetMessage("captchaVerificationFailed", nil), "")
		return
	}

	parsedURL, err := tool.EncodeURI(req.URL)
	if err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("invalidUrl", nil), "URL Parsed Failed")
		log.ErrorPrint("URL Parsed failed: %s", err)
		return
	}

	req.URL = parsedURL.String()
	req.MEMO = url.QueryEscape(req.MEMO)

	if req.EXPIRE != 0 && req.EXPIRE < now {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("illegalExpirationTime", nil), "")
		return
	}

	if !setting.Cfg.AllowAllProtocol && !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("invalidUrl", nil), "Not Allowed Protocol")
		log.WarnPrint("Illegal URL: %s", req.URL)
		return
	}

	link := shorten.GenerateShortenLink(req)

	table := db.SetModel(setting.Cfg.MongoDB.Database, "links")
	err = table.InsertOne(link, link.ShortHash, false)
	if err != nil {
		model.FailureResponse(c, http.StatusInternalServerError, http.StatusInternalServerError, localizer.GetMessage("databaseOperationFailed", nil), "")
		return
	}

	log.DebugPrint("SrcLink: %s, GenerateShortenLink: %s/s/%s", req.URL, setting.Cfg.HTTP.BasePath, link.ShortHash)

	data := map[string]interface{}{
		"hash":  link.ShortHash,
		"token": link.Token,
	}

	model.SuccessResponse(c, data)
}

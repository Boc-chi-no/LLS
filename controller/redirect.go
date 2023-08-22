package controller

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"linkshortener/db"
	"linkshortener/i18n"
	"linkshortener/lib/ip2location"
	"linkshortener/lib/tool"
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

	if err := c.ShouldBindUri(&req); err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("deserializationFailed", nil), "")
		log.ErrorPrint("Deserialization failed: %s", err)
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("deserializationFailed", nil), "")
		log.ErrorPrint("Deserialization failed: %s", err)
		return
	}

	var res []model.Link
	table := db.SetModel(setting.Cfg.MongoDB.Database, "links")
	_ = table.Find(bson.D{{Key: "_id", Value: req.Hash}, {Key: "delete", Value: false}}, &res)

	if res != nil && len(res) == 1 {
		link := res[0]
		reqPassword := ""
		if link.Password != "" {
			if req.Password != "" {
				passwordHash := sha1.Sum([]byte(tool.ConcatStrings(link.ShortHash, req.Password, tool.Uint32ToBase62String(setting.Cfg.Seed))))
				reqPassword = hex.EncodeToString(passwordHash[:])
			}

			if link.Password != reqPassword {
				if reqPassword != "" {
					log.InfoPrint("password error: %s", req.Hash)
				}
				if req.Detect {
					model.FailureResponse(c, http.StatusUnauthorized, http.StatusUnauthorized, localizer.GetMessage("linkPasswordError", nil), "")
					return
				} else {
					c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings("/#/PasswordRedirect/", req.Hash))
					return
				}
			}
		} else if req.Soft {
			c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings("/#/SoftRedirect/", req.Hash))
			return
		}
		go accessLogWorker(c.ClientIP(), req.Hash, c.Request.Header, time.Now().Unix())
		if req.Detect {
			log.DebugPrint("DetectLink: %s", link.URL)
			data := map[string]interface{}{
				"hash": link.ShortHash,
				"url":  link.URL,
			}
			model.SuccessResponse(c, data)
		} else {
			log.DebugPrint("RedirectLink: %s", link.URL)
			c.Redirect(http.StatusTemporaryRedirect, link.URL)
		}
	} else {
		model.FailureResponse(c, http.StatusNotFound, http.StatusNotFound, localizer.GetMessage("noLinkFound", nil), "")
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

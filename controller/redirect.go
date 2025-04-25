package controller

import (
	"crypto/sha256"
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
	"strings"
	"time"
)

// Redirect This method performs the redirection of the shortened link.
// Usage:
// Just hit {BasePath}/s/4nGHqG ( use the generated hash )
func Redirect(c *gin.Context) {
	req := model.RedirectLinkReq{}
	localizer := i18n.GetLocalizer(c)
	now := time.Now().Unix()

	if err := c.ShouldBindUri(&req); err != nil {
		log.ErrorPrint("Deserialization failed: %s", err)
		c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings(setting.Cfg.HTTP.SoftRedirectBasePath, "/#/Error/DeserializationFailed"))
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		log.ErrorPrint("Deserialization failed: %s", err)
		c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings(setting.Cfg.HTTP.SoftRedirectBasePath, "/#/Error/DeserializationFailed"))
		return
	}

	if req.Hash == "ping" {
		Ping(c)
		return
	}

	var res []model.Link
	table := db.SetModel(setting.Cfg.MongoDB.Database, "links")
	_ = table.Find(bson.D{{Key: "_id", Value: req.Hash}, {Key: "delete", Value: false}}, &res, db.Find().SetKey(req.Hash))

	if res != nil && len(res) == 1 {
		link := res[0]
		reqPassword := ""
		if link.Expire != 0 && now > link.Expire {
			if req.Detect {
				model.FailureResponse(c, http.StatusNotFound, http.StatusNotFound, localizer.GetMessage("linkExpire", nil), "")
				return
			} else {
				log.DebugPrint("Link Expire: %s", req.Hash)
				c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings(setting.Cfg.HTTP.SoftRedirectBasePath, "/#/Error/LinkExpire"))
				return
			}
		}
		if link.Password != "" {
			if req.Password != "" {
				passwordHash := sha256.Sum256([]byte(tool.ConcatStrings(link.ShortHash, req.Password, tool.Uint32ToBase62String(setting.Cfg.Seed))))
				reqPassword = hex.EncodeToString(passwordHash[:])
			}

			if link.Password != reqPassword {
				if reqPassword != "" {
					log.DebugPrint("password error: %s", req.Hash)
				}
				if req.Detect {
					model.FailureResponse(c, http.StatusUnauthorized, http.StatusUnauthorized, localizer.GetMessage("linkPasswordError", nil), "")
					return
				} else {
					c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings(setting.Cfg.HTTP.SoftRedirectBasePath, "/#/PasswordRedirect/", req.Hash))
					return
				}
			}
		} else if req.Soft {
			if req.Detect {
				model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("detectAndSoftMutuallyExclusive", nil), "")
				return
			} else {
				c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings(setting.Cfg.HTTP.SoftRedirectBasePath, "/#/SoftRedirect/", req.Hash))
				return
			}
		}
		go accessLogWorker(c.ClientIP(), req.Hash, c.Request.Header, time.Now().Unix())
		if req.Detect {
			log.DebugPrint("DetectLink: %s", link.URL)
			data := map[string]interface{}{
				"hash":   link.ShortHash,
				"url":    link.URL,
				"expire": link.Expire,
				"memo":   link.Memo,
			}
			model.SuccessResponse(c, data)
		} else {
			log.DebugPrint("RedirectLink: %s", link.URL)
			if !strings.HasPrefix(link.URL, "http://") && !strings.HasPrefix(link.URL, "https://") {
				c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings(setting.Cfg.HTTP.SoftRedirectBasePath, "/#/SoftRedirect/", req.Hash))
			} else {
				c.Redirect(http.StatusTemporaryRedirect, link.URL)
			}
		}
	} else {
		if req.Detect {
			model.FailureResponse(c, http.StatusNotFound, http.StatusNotFound, localizer.GetMessage("noLinkFound", nil), "")
			return
		} else {
			c.Redirect(http.StatusTemporaryRedirect, tool.ConcatStrings(setting.Cfg.HTTP.SoftRedirectBasePath, "/#/Error/NotFound"))
			return
		}
	}
}

func accessLogWorker(ip string, hash string, header http.Header, nowTime int64) {
	location := ip2location.Find(ip)
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
	err := table.InsertOne(linkInfo, hash, true)
	if err != nil {
		log.WarnPrint("Failed to write access log to database!")
	}
}

package controller

import (
	"linkshortener/db"
	"linkshortener/i18n"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"math"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// StatsLink This method provides statistics info for redirections
// Usage:
// {BasePath}/api/stats_link
func StatsLink(c *gin.Context) {
	var req model.ManageLinkReq
	localizer := i18n.GetLocalizer(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("deserializationFailed", nil), err.Error())
		log.ErrorPrint("Deserialization failed: %s", err)
		return
	}

	if req.Page <= 0 || req.Size <= 0 || req.Size > 100 {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("invalidPaginationParameter", nil), "")
		return
	}

	// Initialize session object
	session := sessions.Default(c)
	sessionCaptcha := tool.SafeSessionGet(session, "captcha")

	if sessionCaptcha != req.CAPTCHA {
		session.Delete("captcha")
		_ = session.Save()
		model.FailureResponse(c, http.StatusForbidden, http.StatusForbidden, localizer.GetMessage("captchaVerificationFailed", nil), "")
		return
	}

	var res []model.Link
	table := db.SetModel(setting.Cfg.MongoDB.Database, "links")
	_ = table.Find(bson.D{{Key: "_id", Value: req.Hash}}, &res, db.Find().SetKey(req.Hash))

	if res != nil && len(res) > 0 {
		if res[0].Token != req.Token {
			session.Delete("captcha")
			_ = session.Save()
			model.FailureResponse(c, http.StatusForbidden, http.StatusForbidden, localizer.GetMessage("passwordVerificationFailed", nil), "")
			return
		}

		var statsRes []model.LinkInfo
		statsTable := db.SetModel(setting.Cfg.MongoDB.Database, "link_access")

		offset := (req.Page - 1) * req.Size
		totalCount, _ := statsTable.CountDocuments(bson.D{{Key: "hash", Value: req.Hash}}, db.Find().SetKey(req.Hash))
		totalPages := int64(math.Ceil(float64(totalCount) / float64(req.Size)))

		if totalCount > 0 && req.Page <= totalPages {
			_ = statsTable.Find(bson.D{{Key: "hash", Value: req.Hash}}, &statsRes, db.Find().SetSkip(offset).SetLimit(req.Size).SetKey(req.Hash).SetPrefixScans(true))

			data := map[string]interface{}{
				"current": req.Page,
				"size":    req.Size,

				"pages":   totalPages,
				"total":   totalCount,
				"records": statsRes,
			}

			model.SuccessResponse(c, data)
		} else {
			data := map[string]interface{}{
				"current": req.Page,
				"size":    req.Size,
				"pages":   0,
				"total":   0,
				"records": []string{},
			}
			model.SuccessResponse(c, data)
		}

	} else {
		session.Delete("captcha")
		_ = session.Save()
		model.FailureResponse(c, http.StatusNotFound, http.StatusNotFound, localizer.GetMessage("noLinkFound", nil), "")
	}
}

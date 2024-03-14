package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"linkshortener/db"
	"linkshortener/i18n"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
)

// DeleteLink This method deletes the redirection
// Usage:
// Send a http POST call to
// {BasePath}/api/delete_link
func DeleteLink(c *gin.Context) {
	var req model.ManageLinkReq
	localizer := i18n.GetLocalizer(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		model.FailureResponse(c, http.StatusBadRequest, http.StatusBadRequest, localizer.GetMessage("deserializationFailed", nil), err.Error())
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

	var res []model.Link
	table := db.SetModel(setting.Cfg.MongoDB.Database, "links")
	_ = table.Find(bson.D{{Key: "_id", Value: req.Hash}, {Key: "delete", Value: false}}, &res, db.Find().SetKey(req.Hash))

	if res != nil && len(res) > 0 {
		if res[0].Token != req.Token {
			model.FailureResponse(c, http.StatusForbidden, http.StatusForbidden, localizer.GetMessage("passwordVerificationFailed", nil), "")
			return
		}
		err := table.UpdateByID(req.Hash, bson.M{
			"$set": bson.M{
				"delete": true,
			},
		})

		if err != nil {
			model.FailureResponse(c, http.StatusInternalServerError, http.StatusInternalServerError, localizer.GetMessage("databaseOperationFailed", nil), "")
			return
		}
		model.SuccessResponse(c, nil)
	} else {
		model.FailureResponse(c, http.StatusNotFound, http.StatusNotFound, localizer.GetMessage("noLinkFound", nil), "")
	}

}

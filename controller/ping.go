package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"linkshortener/db"
	"linkshortener/i18n"
	"linkshortener/lib/tool"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
	"time"
)

// Ping Verify that the server is available.
// Usage:
// Just hit {BasePath}/ping
func Ping(c *gin.Context) {
	localizer := i18n.GetLocalizer(c)

	var res []model.Link
	table := db.SetModel(setting.Cfg.MongoDB.Database, "links")
	_ = table.Find(bson.D{{Key: "_id", Value: "000000"}}, &res, db.Find().SetKey("000000"))

	if res != nil && len(res) == 1 {
		model.SuccessResponse(c, map[string]interface{}{
			"msg": "pong",
		})
		return
	} else {
		var link model.Link
		link.ShortHash = "000000"
		link.URL = "http://www.example.com/"
		link.Created = time.Now().Unix()
		link.Token, _ = tool.GetToken(16)
		link.Memo = "Test Hash"
		link.Delete = false
		err := table.InsertOne(link, link.ShortHash, false)
		if err != nil {
			model.FailureResponse(c, http.StatusInternalServerError, http.StatusInternalServerError, localizer.GetMessage("databaseOperationFailed", nil), "")
			return
		} else {
			model.SuccessResponse(c, map[string]interface{}{
				"msg": "pong",
			})
			return
		}
	}
}

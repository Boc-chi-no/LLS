package controller

import (
	"linkshortener/db"
	"linkshortener/i18n"
	"linkshortener/lib/tool"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Ping Verify that the server is available.
// Usage:
// Just hit {BasePath}/ping
func Ping(c *gin.Context) {
	localizer := i18n.GetLocalizer(c)

	var res []model.Link
	table := db.SetModel(setting.Cfg.DB.Database, "links")
	_ = table.Find(bson.D{{Key: "_id", Value: "000000"}}, &res, db.Find().SetKey("000000"))

	if res != nil && len(res) == 1 {
		model.SuccessResponse(c, map[string]interface{}{
			"msg": "pong",
		})
		return
	} else {
		var link model.Link
		link.ShortHash = "000000"
		link.URL = "https://www.example.com/"
		link.Created = time.Now().Unix()
		link.Token, _ = tool.GetToken(16)
		link.Memo = "Test Hash"
		link.Delete = false
		_, err := table.InsertOne(link, false)
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

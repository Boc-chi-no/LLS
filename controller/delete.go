package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"linkshortener/db"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
)

// DeleteLink This method deletes the redirection
// Usage:
// Send a http POST call to
// http://localhost:8040/api/delete_link
func DeleteLink(c *gin.Context) {
	var req model.ManageLinkReq

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

	var res []model.Link
	table := db.NewModel(setting.Cfg.MongoDB.Database, "links")
	table.Find(bson.D{{Key: "_id", Value: req.Hash}, {Key: "delete", Value: false}}, &res)

	if res != nil && len(res) > 0 {
		if res[0].Token != req.Token {
			model.FailureResponse(c, http.StatusForbidden, http.StatusForbidden, "链接密码检验失败!", "")
			return
		}
		table.UpdateByID(req.Hash, bson.M{
			"$set": bson.M{
				"delete": true,
			},
		})
		model.SuccessResponse(c, nil)
	} else {
		model.FailureResponse(c, 404, 404, "未找到查询的链接!", "")
	}

}

package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"linkshortener/db"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"math"
	"net/http"
)

// StatsLink This method provides statistics info for redirections
// Usage:
// http://localhost:8040/api/stats_link
func StatsLink(c *gin.Context) {
	var req model.ManageLinkReq

	if err := c.BindJSON(&req); err != nil {
		model.FailureResponse(c,http.StatusBadRequest,http.StatusBadRequest,"序列化失败","")
		log.ErrorPrint("Deserialization failed: %s",err)
		return
	}

	if req.Page <= 0||req.Size<= 0||req.Size> 100{
		model.FailureResponse(c,http.StatusBadRequest,http.StatusBadRequest,"错误的分页参数","")
		return
	}

	// 初始化session对象
	session := sessions.Default(c)
	sessionCaptcha := session.Get("captcha")
	session.Delete("captcha")
	_ = session.Save()

	if sessionCaptcha != req.CAPTCHA {
		model.FailureResponse(c,http.StatusForbidden,http.StatusForbidden,"验证码检验失败!","")
		return
	}

	var res []model.Link
	table:=db.NewModel(setting.Cfg.MongoDB.Database, "links")
	table.Find(bson.D{{Key:"_id",Value: req.Hash}},&res)

	if res != nil&&len(res) > 0 {
		if res[0].Token != req.Token {
			model.FailureResponse(c,http.StatusForbidden,http.StatusForbidden,"链接密码检验失败!","")
			return
		}

		var statsRes []model.LinkInfo
		statsTable:=db.NewModel(setting.Cfg.MongoDB.Database, "link_access")

		offset := (req.Page - 1) * req.Size
		totalCount := statsTable.CountDocuments(bson.D{{Key:"hash",Value: req.Hash}})
		totalPages := int64(math.Ceil(float64(totalCount) / float64(req.Size)))

		if totalCount>0&&req.Page<=totalPages{
			statsTable.Find(bson.D{{Key:"hash",Value: req.Hash}},&statsRes,options.Find().SetSkip(offset).SetLimit(req.Size))

			data := map[string]interface{}{
				"current": req.Page,
				"size": req.Size,

				"pages": totalPages,
				"total": totalCount,
				"records":statsRes,
			}

			model.SuccessResponse(c,data)
		}else {
			data := map[string]interface{}{
				"current": req.Page,
				"size": req.Size,
				"pages": 0,
				"total": 0,
				"records":[]string{},
			}
			model.SuccessResponse(c,data)
		}

	} else {
		model.FailureResponse(c,404,404,"未找到查询的链接!","")
	}
}

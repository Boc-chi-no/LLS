package controller

import (
	"bytes"
	"encoding/base64"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"image/png"
	"linkshortener/lib/captcha"
	"linkshortener/log"
	"linkshortener/model"
	"net/http"
	"strings"
)

// Captcha This method generate captcha.
//
// Usage:
// generate captcha, just http GET to http://localhost:8040/api/captcha
func Captcha(c *gin.Context) {
	// 初始化session对象
	session := sessions.Default(c)
	cp := captcha.NewCaptcha(120, 40, 4)
	cp.SetMode(1) // 设置为数学公式模式
	code, img := cp.OutPut()

	// 设置session数据
	session.Set("captcha", code)
	session.Options(sessions.Options{
		MaxAge:  300,
	})
	err := session.Save()
	if err != nil {
		model.FailureResponse(c,http.StatusInternalServerError,http.StatusInternalServerError,"保存session失败","")
		log.ErrorPrint("session save failed: %s",err)
		return
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		model.FailureResponse(c,http.StatusInternalServerError,http.StatusInternalServerError,"图片生成失败","")
		log.ErrorPrint("captcha image generation failed: %s",err)
		return
	}

	var imgBase64Str strings.Builder
	imgBase64Str.WriteString("data:image/png;base64,")
	imgBase64Str.WriteString(base64.StdEncoding.EncodeToString(buf.Bytes()))

	data := map[string]interface{}{
		"pic": imgBase64Str.String(),
	}

	c.Header("Cache-Control", "no-store")
	model.SuccessResponse(c,data)
}

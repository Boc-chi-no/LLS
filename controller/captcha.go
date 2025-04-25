package controller

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"linkshortener/i18n"
	"linkshortener/lib/captcha"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/model"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Captcha This method generate captcha.
//
// Usage:
// generate captcha, just http GET to {BasePath}/api/captcha
func Captcha(c *gin.Context) {
	localizer := i18n.GetLocalizer(c)
	// Initialize session object
	session := sessions.Default(c)
	cp := captcha.NewCaptcha(120, 40, 4)
	cp.SetMode(1) // Set to maths mode
	code, img := cp.OutPut()

	// Setting session data
	_ = tool.SafeSessionGet(session, "captcha")
	session.Set("captcha", code)
	session.Options(sessions.Options{
		MaxAge: 300,
	})
	err := session.Save()
	if err != nil {
		model.FailureResponse(c, http.StatusInternalServerError, http.StatusInternalServerError, localizer.GetMessage("saveSessionFailed", nil), "")
		log.ErrorPrint("session save failed: %s", err)
		return
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		model.FailureResponse(c, http.StatusInternalServerError, http.StatusInternalServerError, localizer.GetMessage("imageGenerationFailed", nil), "")
		log.ErrorPrint("captcha image generation failed: %s", err)
		return
	}

	var imgBase64Str strings.Builder
	imgBase64Str.WriteString("data:image/png;base64,")
	imgBase64Str.WriteString(base64.StdEncoding.EncodeToString(buf.Bytes()))

	data := map[string]interface{}{
		"pic": imgBase64Str.String(),
	}

	c.Header("Cache-Control", "no-store")
	model.SuccessResponse(c, data)
}

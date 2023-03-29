package i18n

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/setting"
	"os"
)

var bundle *i18n.Bundle

type Translator struct {
	localizer *i18n.Localizer
}

type ITranslator interface {
	GetMessage(id string, templateData map[string]interface{}) string
}

func bundleMustLoadMessageBytes(buf []byte, tag string) {
	_, err := bundle.ParseMessageFileBytes(buf, tool.ConcatStrings("LLS.", tag, ".json"))
	if err != nil {
		log.PanicPrint("Failed to initialize LanguageFiles[%s]: %s", tag, err)
	}
}

// InitI18n Initialize a global variable `bundle` and set the default language to English.
// Load language resource files during initialization.
func InitI18n(jp []byte, cn []byte, us []byte) {
	bundle = i18n.NewBundle(language.AmericanEnglish)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	if setting.Cfg.I18N.AddExtraLanguage {
		bytes, err := os.ReadFile(setting.Cfg.I18N.ExtraLanguageFiles)
		if err != nil {
			log.PanicPrint("Failed to Read ExtraLanguageFiles: %s", err)
		}
		bundleMustLoadMessageBytes(bytes, setting.Cfg.I18N.ExtraLanguageName)
	}
	bundleMustLoadMessageBytes(jp, "ja-JP")
	bundleMustLoadMessageBytes(cn, "zh-CN")
	bundleMustLoadMessageBytes(us, "en-US")
}

// GetLocalizer Retrieve the preferred language based on the 'Language' header of the user's request.
// If not available, use the 'Accept-Language' header to set the language.
// return a Translator object.
func GetLocalizer(c *gin.Context) ITranslator {
	lang := c.GetHeader("Language")
	acceptLanguage := c.GetHeader("Accept-Language")
	return &Translator{
		localizer: i18n.NewLocalizer(bundle, lang, acceptLanguage),
	}
}

// GetMessage retrieves the localized message corresponding to the given messageID,
// and optionally applies template data to the message before returning it.
// If the message is not found or an error occurs during localization, this method
// will panic. It returns the localized message as a string.
func (t Translator) GetMessage(messageID string, templateData map[string]interface{}) string {
	return t.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}

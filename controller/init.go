package controller

import (
	"context"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"linkshortener/fs"
	"linkshortener/i18n"
	"linkshortener/lib/lfs"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"net/http"
	"strings"
	"sync"
	"time"
)

var router *gin.Engine

func ReqLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method

		reqUri := c.Request.RequestURI

		statusCode := c.Writer.Status()

		clientIP := c.ClientIP()

		log.InfoPrint("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func InitRouter() {
	BasePath := strings.TrimPrefix(strings.TrimSuffix(strings.Join(strings.Fields(setting.Cfg.HTTP.BasePath), ""), "/"), "/")
	if BasePath != "" {
		BasePath = tool.ConcatStrings("/", BasePath)
	}
	SoftRedirectBasePath := strings.TrimPrefix(strings.TrimSuffix(strings.Join(strings.Fields(setting.Cfg.HTTP.SoftRedirectBasePath), ""), "/"), "/")
	if SoftRedirectBasePath != "" {
		SoftRedirectBasePath = tool.ConcatStrings("/", SoftRedirectBasePath)
	}

	setting.Cfg.HTTP.BasePath = BasePath
	setting.Cfg.HTTP.SoftRedirectBasePath = SoftRedirectBasePath

	router.GET(tool.ConcatStrings(BasePath, "/ping"), func(c *gin.Context) {
		model.SuccessResponse(c, map[string]interface{}{
			"msg": "pong",
		})
	}) //Service Test Interface

	router.GET(tool.ConcatStrings(BasePath, "/s/:hash"), Redirect) //Short link redirection

	router.GET(tool.ConcatStrings(BasePath, "/api/captcha"), Captcha)             //Generate captcha code
	router.POST(tool.ConcatStrings(BasePath, "/api/generate_link"), GenerateLink) //Create link
	router.POST(tool.ConcatStrings(BasePath, "/api/stats_link"), StatsLink)       //Link statistics
	router.POST(tool.ConcatStrings(BasePath, "/api/delete_link"), DeleteLink)     //Delete link

	if setting.Cfg.HTTP.DisableFilesDirEmbed { //Static files
		if strings.Join(strings.Fields(setting.Cfg.HTTP.FilesDirURI), "") != "" {
			router.NoRoute(gin.WrapH(http.FileServer(
				lfs.LlsFileSystem{
					Fs: http.Dir(setting.Cfg.HTTP.FilesDirURI), //Use of external resources
				},
			)))
		} else {
			log.PanicPrint("STATIC_FILES_DIR_URI not allowed to be empty")
		}
	} else {
		router.NoRoute(gin.WrapH(tool.HTTPAddPrefix("/ui", http.FileServer(
			lfs.LlsFileSystem{
				Fs: fs.StatikFS, //Use of embedded resources
			},
		))))
	}
}

func NewLimiter(reqRate rate.Limit, reqBurst int, reqTimeout time.Duration) gin.HandlerFunc {
	limiters := &sync.Map{}

	return func(c *gin.Context) {
		if c.FullPath() != "" {
			key := c.ClientIP()
			limit, _ := limiters.LoadOrStore(key, rate.NewLimiter(reqRate, reqBurst))

			ctx, cancel := context.WithTimeout(c, reqTimeout)
			defer cancel()

			if err := limit.(*rate.Limiter).Wait(ctx); err != nil {
				localizer := i18n.GetLocalizer(c)
				model.FailureResponse(c, http.StatusTooManyRequests, http.StatusTooManyRequests, localizer.GetMessage("tooManyRequests", nil), "")
			}
		}
		c.Next()
	}
}

func LooseCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func InitController() {
	if setting.Cfg.RunMode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.DefaultWriter = log.NullOut
		gin.DefaultErrorWriter = log.NullOut
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()

	if setting.Cfg.HTTP.LooseCORS {
		router.Use(LooseCORS())
	}

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.DebugPrint("%v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	router.Use(ReqLogger())

	if setting.Cfg.HTTPLimiter.EnableLimiter {
		router.Use(NewLimiter(rate.Limit(setting.Cfg.HTTPLimiter.LimitRate), setting.Cfg.HTTPLimiter.LimitBurst, time.Duration(setting.Cfg.HTTPLimiter.Timeout)*time.Millisecond))
	}

	SessionSecret, _ := tool.GetToken(16)
	if !setting.Cfg.HTTP.RandomSessionSecret {
		SessionSecret = setting.Cfg.HTTP.SessionSecret
	}
	store := memstore.NewStore([]byte(SessionSecret))
	router.Use(sessions.Sessions("session", store))

	if setting.Cfg.RunMode == "dev" {
		pprof.Register(router) //debug
	}
}

func RunServer() {
	log.InfoPrint("Listening and serving HTTP on %s", setting.Cfg.HTTP.Listen)
	err := router.Run(setting.Cfg.HTTP.Listen)
	if err != nil {
		log.PanicPrint("Start Web Server Fail: %s", err)
	}
}

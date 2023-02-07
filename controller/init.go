package controller

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"linkshortener/statikFS"
	"net/http"
	"time"
)

var router *gin.Engine

func ReqLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
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
	router.GET("ping", func(c *gin.Context) {
		model.SuccessResponse(c, map[string]interface{}{
			"msg": "pong",
		})
	}) //服务测试接口

	router.GET("/s/:hash", Redirect) //短链接重定向

	router.GET("/api/captcha", Captcha)            //生成验证码
	router.POST("/api/generate_link", InsertLink)  //创建链接
	router.POST("/api/stats_link", StatsLink)      //链接统计
	router.POST("/api/delete_link", DeleteLink)    //删除链接

	if setting.Cfg.HTTP.FilesEmbed { //静态文件
		router.NoRoute(gin.WrapH(http.FileServer(statikFS.StatikFS))) //使用内嵌资源
	} else {
		router.NoRoute(gin.WrapH(http.FileServer(http.Dir(setting.Cfg.HTTP.FilesURI)))) //使用外部资源
	}
}

func InitController() {

	if setting.Cfg.RunMode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.DebugPrint("%v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	router.Use(ReqLogger())

	if setting.Cfg.HTTPLimiter.EnableLimiter {
		router.Use(tool.NewLimiter(rate.Limit(setting.Cfg.HTTPLimiter.LimitRate), setting.Cfg.HTTPLimiter.LimitBurst, time.Duration(setting.Cfg.HTTPLimiter.Timeout)*time.Millisecond))
	}

	SessionSecret := tool.GetToken(16)
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
	log.InfoPrint("Listening and serving HTTP on %s",setting.Cfg.HTTP.Listen)
	err := router.Run(setting.Cfg.HTTP.Listen)
	if err != nil {
		log.PanicPrint("Start Web Server Fail: %s",err)
	}
}
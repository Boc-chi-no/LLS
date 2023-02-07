package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"linkshortener/model"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func Time() int {
	cur := time.Now()
	timestamp := cur.UnixNano() / 1000000
	return int(timestamp / 1000)
}

func Now() string {
	tm := time.Unix(int64(Time()), 0)
	return tm.Format("2006-01-02 15:04:05")
}

func NowDay() string {
	tm := time.Unix(int64(Time()), 0)
	return tm.Format("20060102")
}

//StringToObject json字符串转对象
func StringToObject(str string, data interface{}) bool {
	js := json.NewDecoder(bytes.NewReader([]byte(str)))
	js.UseNumber()
	err := js.Decode(data)
	if err == nil {
		return true
	}
	return false
}

//Mkdir 创建目录
func Mkdir(path string) {
	_, e1 := os.Stat(path)
	if e1 != nil && !os.IsExist(e1) {
		if path != "" {
			os.MkdirAll(path, 0777)
		}
	}
}

//FileExist 判断文件是否存在，存在返回true 不存在返回false
func FileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	} else {
		return false
	}
}

func GetToken(length int) string {
	sb := strings.Builder{}
	sb.Grow(length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}


func NewLimiter(reqRate rate.Limit, reqBurst int, reqTimeout time.Duration) gin.HandlerFunc {
	limiters := &sync.Map{}

	return func(c *gin.Context) {
		if c.FullPath() != ""{
			key := c.ClientIP()
			limit, _ := limiters.LoadOrStore(key, rate.NewLimiter(reqRate, reqBurst))

			ctx, cancel := context.WithTimeout(c, reqTimeout)
			defer cancel()

			if err := limit.(*rate.Limiter).Wait(ctx); err != nil {
				model.FailureResponse(c,http.StatusTooManyRequests,http.StatusTooManyRequests, "请求过于频繁","")
			}
		}
		c.Next()
	}
}
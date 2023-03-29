package tool

import (
	"math/rand"
	"os"
	"strings"
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

func ConcatStrings(str ...string) string {
	var sb strings.Builder
	for _, s := range str {
		sb.WriteString(s)
	}
	return sb.String()
}

// StringToObject is a utility function that takes a JSON string and decodes it into the given data object.
// The function returns a boolean value indicating whether the decoding was successful or not.
//func StringToObject(str string, data interface{}) bool {
//	js := json.NewDecoder(bytes.NewReader([]byte(str)))
//	js.UseNumber()
//	err := js.Decode(data)
//	if err == nil {
//		return true
//	}
//	return false
//}

// Mkdir creates a directory at the given path if it does not already exist.
// Returns an error if the directory could not be created or if an error occurred during stat.
func Mkdir(path string) error {
	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		if path != "" {
			err = os.MkdirAll(path, 0777)
		}
	}
	return err
}

// FileExist checks whether the given file exists or not
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

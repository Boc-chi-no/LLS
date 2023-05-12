package tool

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var Base62Map = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
	"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
	"u", "v", "w", "x", "y", "z", "A", "B", "C", "D",
	"E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
	"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
	"Y", "Z",
}

var GlobalCounter int64

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

func Uint32ToBase62String(uitNum uint32) string {
	num := int(uitNum)
	var remainder int
	var base62Str strings.Builder

	if num == 0 {
		return "0"
	} else {
		for num != 0 {
			remainder = num % 62
			if 76 > remainder && remainder > 9 {
				base62Str.WriteString(Base62Map[remainder])
			} else {
				base62Str.WriteString(strconv.Itoa(remainder))
			}
			num = num / 62
		}
		return base62Str.String()
	}
}

func GetToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length %d, must be greater than zero", length)
	}

	bytes := make([]byte, length*4)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %s", err)
	}

	var token strings.Builder
	token.Grow(length)
	for i := 0; i < length; i++ {
		uint32Slice := binary.LittleEndian.Uint32(bytes[i*4 : i*4+4])
		token.WriteString(Base62Map[int(uint32Slice%62)])
	}
	return token.String(), nil
}

func GlobalCounterSafeAdd(delta int64) int64 {
	return atomic.AddInt64(&GlobalCounter, delta)
}

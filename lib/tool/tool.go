package tool

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"os"
	"reflect"
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

var GlobalCounter atomic.Uint64

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
			err = os.MkdirAll(path, 0700)
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

func GlobalCounterSafeAdd(delta uint64) uint64 {
	return GlobalCounter.Add(delta)
}

func HTTPAddPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = ConcatStrings(prefix, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

// IsDataMatchingFilter checks if data matches the filter conditions
func IsDataMatchingFilter(data, filter map[string]interface{}) bool {
	// Iterate over key-value pairs in the filter
	for key, filterValue := range filter {
		// Check if the data contains the key from the filter
		dataValue, keyExists := data[key]
		if !keyExists {
			// Data does not contain the key from the filter, not a match
			return false
		}

		// Check if the value in the data equals the value in the filter
		if dataValue != filterValue {
			// Values do not match, not a match
			return false
		}
	}

	// Data matches all filter conditions
	return true
}

func processStructFields(val reflect.Value, typ reflect.Type, resultMap map[string]interface{}) error {
	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)

		// Get the "bson" tag; if not present, use the field name as the tag
		tag := fieldType.Tag.Get("bson")
		if tag == "" {
			tag = fieldType.Name
		}

		// Check if the field is exportable (i.e., starts with an uppercase letter)
		if fieldType.PkgPath != "" {
			return fmt.Errorf("MarshalJsonByBson: field %s is not exportable", fieldType.Name)
		}

		// If the field is a struct, recursively process its fields
		if fieldValue.Kind() == reflect.Struct {
			nestedMap := make(map[string]interface{})
			if err := processStructFields(fieldValue, fieldType.Type, nestedMap); err != nil {
				return err
			}
			resultMap[tag] = nestedMap
		} else {
			// Add non-struct field to the JSON map
			resultMap[tag] = fieldValue.Interface()
		}
	}
	return nil
}

// unmarshalStructFields recursively processes struct fields and handles nested structs.
// It takes a reflect.Value representing a struct, a map of JSON values, and returns a map of processed values.
func unmarshalStructFields(structValue reflect.Value, resultMap map[string]interface{}) (map[string]interface{}, error) {
	var err error
	processedMap := make(map[string]interface{})
	if structValue.Kind() != reflect.Struct {
		return processedMap, nil
	}

	for i := 0; i < structValue.NumField(); i++ {
		fieldReflectValue := structValue.Field(i)
		fieldType := structValue.Type().Field(i)

		// Get the "bson" tag; if not present, use the field name as the tag
		tag := fieldType.Tag.Get("bson")
		if tag == "" {
			tag = fieldType.Name
		}

		// Check if the tag exists in the JSON map
		if jsonValue, ok := resultMap[tag]; ok {
			if fieldReflectValue.Kind() == reflect.Struct {
				nestedMap, nestedMapOk := jsonValue.(map[string]interface{})
				if nestedMapOk {
					nestedStructMap, nestedStructErr := unmarshalStructFields(fieldReflectValue, nestedMap)
					if nestedStructErr != nil {
						return nil, nestedStructErr
					}

					for key, value := range nestedStructMap {
						processedMap[key] = value
					}
				}
			} else {
				processedMap[fieldType.Name] = jsonValue
			}
		}
	}

	return processedMap, err
}

// MarshalJsonByBson serializes a struct into JSON data based on "bson" tags.
// It takes a struct and returns the JSON representation as a byte slice.
func MarshalJsonByBson(i interface{}) ([]byte, error) {
	val := reflect.ValueOf(i)
	typ := val.Type()
	jsonMap := make(map[string]interface{})
	// Ensure that i is a valid struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("results argument must be a pointer to a slice, but was a pointer to %s", val.Kind())
	}

	// If i is a pointer, dereference it to get the underlying struct value
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Recursively process struct fields
	if err := processStructFields(val, typ, jsonMap); err != nil {
		return nil, err
	}

	// Marshal the JSON map into a byte slice
	data, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, fmt.Errorf("MarshalJsonByBson: failed to marshal JSON: %v", err)
	}

	return data, nil
}

// UnmarshalJsonByBson deserializes JSON data into a target slice of structs based on "bson" tags.
// It takes a byte slice containing JSON data and a pointer to a slice of the target struct for deserialization.
func UnmarshalJsonByBson(data []byte, i interface{}) error {
	resultsVal := reflect.ValueOf(i)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}
	sliceVal := resultsVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind())
	}
	sliceType := sliceVal.Type().Elem()

	// Unmarshal the JSON data into a slice of maps
	var jsonSlice []map[string]interface{}
	newJsonSlice := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(data, &jsonSlice); err != nil {
		return fmt.Errorf("UnmarshalJsonByBson: failed to unmarshal JSON data: %v", err)
	}

	// Iterate through the JSON slice and create structs
	for _, jsonMap := range jsonSlice {
		structValue := reflect.New(sliceType).Elem()
		newJson, err := unmarshalStructFields(structValue, jsonMap)
		if err != nil {
			return err
		}

		newJsonSlice = append(newJsonSlice, newJson)
	}

	mSliceJson, _ := json.Marshal(newJsonSlice)

	return json.Unmarshal(mSliceJson, i)
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

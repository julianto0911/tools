package tools

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/jinzhu/copier"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func PadLeft(val, length int) string {
	return fmt.Sprintf("%0*d", length, val)
}

func ToJSONString(val interface{}) string {
	temp, _ := json.Marshal(val)
	return string(temp)
}

func ToBytes(val interface{}) []byte {
	temp, _ := json.Marshal(val)
	return temp
}

func StructToUrlValue(data interface{}) (url.Values, error) {
	return query.Values(data)
}

func Substring(str string, start, end int) string {
	return strings.TrimSpace(str[start : end+1])
}

func StringToFloat64(val string) float64 {
	t, _ := strconv.ParseFloat(strings.TrimSpace(val), 64)
	return t
}

func Float64ToString(val float64) string {
	return fmt.Sprintf("%v", val)
}

func StringToInt64(val string) int64 {
	t, _ := strconv.ParseInt(val, 10, 64)
	return t
}

func StringToInt8(val string) int8 {
	t, _ := strconv.ParseInt(val, 10, 8)
	return int8(t)
}

func StringToInt16(val string) int16 {
	t, _ := strconv.ParseInt(val, 10, 16)
	return int16(t)
}

func StringToInt32(val string) int32 {
	t, _ := strconv.ParseInt(val, 10, 32)
	return int32(t)
}

func StringToInt(val string) int {
	t, _ := strconv.Atoi(val)
	return t
}

func ConvertRawAmount(v string) float64 {
	w := strings.ReplaceAll(v, ".", "")
	resp, _ := strconv.ParseFloat(w, 64)
	return resp
}

func NumericStringToInt64(v string) int64 {
	w := strings.ReplaceAll(v, ".", "")
	return StringToInt64(w)
}

func NumericStringToFloat64(v string) float64 {
	w := strings.ReplaceAll(v, ".", "")
	return StringToFloat64(w)
}

func ToStrThousands(v float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%d", int64(v))
}

func ToStrThousandsI(v int64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%d", (v))
}

func IntToString(v int) string {
	return strconv.Itoa(v)
}

func Int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

func DayOfWeek(v time.Time) int {
	return int(v.Weekday())
}

// Generic converter
func ToDTO[E any, D any](entity E) (*D, error) {
	var dto D
	if err := copier.Copy(&dto, &entity); err != nil {
		return nil, err
	}
	return &dto, nil
}

// Generic slice converter
func ToDTOs[E any, D any](entities []E) ([]D, error) {
	dtos := make([]D, len(entities))
	for i, entity := range entities {
		dto, err := ToDTO[E, D](entity)
		if err != nil {
			return nil, err
		}
		dtos[i] = *dto
	}
	return dtos, nil
}

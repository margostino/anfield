package common

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

func RegisterTime(timeType string, requestId int) {
	fmt.Printf("%s #%d: %s\n", timeType, requestId, time.Now().String())
}

func WaitGroup(delta int) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(delta)
	return &wg
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func IsError(e error) bool {
	return e != nil
}

func UTC(time time.Time) string {
	formattedTime := time.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s (UTC)", formattedTime)
}

func InSlice(value, slice interface{}) bool {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice, reflect.Ptr:
		values := reflect.Indirect(reflect.ValueOf(slice))
		if values.Len() == 0 {
			return false
		}

		val := reflect.Indirect(reflect.ValueOf(value))

		if val.Kind() != values.Index(0).Kind() {
			return false
		}

		for i := 0; i < values.Len(); i++ {
			if reflect.DeepEqual(values.Index(i).Interface(), val.Interface()) {
				return true
			}
		}
	}
	return false
}

func Hash(value string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(value))
	return hash.Sum32()
}

func NormalizeDay(day string) string {
	if len(day) == 1 {
		return "0" + day
	}
	return day
}

func NormalizeMonth(name string) string {
	switch strings.ToLower(name) {
	case "january":
		return "01"
	case "february":
		return "02"
	case "march":
		return "03"
	case "april":
		return "04"
	case "may":
		return "05"
	case "june":
		return "06"
	case "july":
		return "07"
	case "august":
		return "08"
	case "september":
		return "09"
	case "october":
		return "10"
	case "november":
		return "11"
	case "december":
		return "12"
	default:
		return "-1"
	}
}

func IsTimeCounter(value string) bool {
	match, _ := regexp.MatchString("([0-9]?'|[0-9]{2}'|[0-9]{2}\\+[0-9]+'|HT)$", value)
	return match
}

func IsFormationNumber(value string) bool {
	match, _ := regexp.MatchString("^\\b[0-9]{1,2}\\b$", value)
	return match
}

func Even(number int) bool {
	return number%2 == 0
}

func Remove(slice []string, element string) []string {
	var newSlice []string
	for _, value := range slice {
		if element != value {
			newSlice = append(newSlice, value)
		}
	}

	if newSlice != nil {
		return newSlice
	}

	return slice
}

func GenerateEventID(url string) string {
	home, away, identifier := ExtractTeamsFrom(url)
	return fmt.Sprintf("%s-%s-%s", home, away, identifier)
}

func ExtractTeamsFrom(url string) (string, string, string) {
	var home, away, identifier string
	substringsPath := strings.Split(url, "/")

	if len(substringsPath) > 8 {
		identifier = substringsPath[8]
	}

	if len(substringsPath) > 7 {
		h2h := strings.Split(substringsPath[7], "vs")

		if len(h2h) > 1 {
			home = strings.ReplaceAll(strings.Trim(h2h[0], "-"), "-", "_")
			away = strings.ReplaceAll(strings.Trim(h2h[1], "-"), "-", "_")
		}
	}

	if identifier == "" {
		home = "invalid_identifier"
	}

	if home == "" {
		home = "invalid_home"
	}

	if away == "" {
		away = "invalid_away"
	}

	return home, away, identifier
}

func HashFrom(seed string) string {
	hash := sha1.New()
	io.WriteString(hash, seed)
	return hex.EncodeToString(hash.Sum(nil))
}

func Now() time.Time {
	return time.Now().UTC()
}

func Mask(s string) string {
	prefix := s[:4]
	suffix := s[len(s)-4:]
	return prefix + "****" + suffix
}

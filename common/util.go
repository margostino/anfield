package common

import (
	"fmt"
	"hash/fnv"
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
		panic(e)
	}
}

func Contains(list []uint32, element uint32) bool {
	for _, value := range list {
		if value == element {
			return true
		}
	}
	return false
}

func Hash(value string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(value))
	return hash.Sum32()
}

func Reverse(list *[]string) {
	for i := 0; i < len(*list)/2; i++ {
		j := len(*list) - i - 1
		(*list)[i], (*list)[j] = (*list)[j], (*list)[i]
	}
}

func NormalizeMonth(name string) string {
	switch name {
	case "January":
		return "01"
	case "February":
		return "02"
	case "March":
		return "03"
	case "April":
		return "04"
	case "May":
		return "05"
	case "June":
		return "06"
	case "July":
		return "07"
	case "August":
		return "08"
	case "September":
		return "09"
	case "October":
		return "10"
	case "November":
		return "11"
	case "December":
		return "12"
	default:
		return "-1"
	}
}

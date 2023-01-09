package cache

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
)

func parseSize(size string) (int64, string) {
	//默认大小为100MB
	re, _ := regexp.Compile("[0-9]+")
	unit := string(re.ReplaceAll([]byte(size), []byte("")))                //将数字部分替换为"", B KB MB
	num, _ := strconv.ParseInt(strings.Replace(size, unit, "", 1), 10, 64) //将字母替换为""
	var byteNum int64
	switch unit {
	case "B":
		byteNum = num * B
	case "KB":
		byteNum = num * KB
	case "MB":
		byteNum = num * MB
	case "GB":
		byteNum = num * GB
	case "TB":
		byteNum = num * TB
	case "PB":
		byteNum = num * PB
	default:
		num = 0
	}
	if num == 0 {
		num = 100
		byteNum = num * MB
		unit = "MB"
	}
	sizeStr := strconv.FormatInt(num, 10) + unit
	return byteNum, sizeStr
}

func getValSize(val interface{}) int64 {
	bytes, _ := json.Marshal(val)
	return int64(len(bytes))
}

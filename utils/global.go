package utils

import (
	"regexp"
	"strings"
	"strconv"
)

func HasValue(s string) bool {
	if len(s) == 0 {
		return false
	}
	return true
}

func IsPhoneNumber(s string) bool {
	matched, _ := regexp.Match(`\d{7,12}`, []byte(s))
	return matched
}

func IsExtensionNumber(s string) bool {
	matched, _ := regexp.Match(`\d{4,12}`, []byte(s))
	return matched
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// camel string, xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data[:])
}

/*
 * 被加数：string
 * 加数：int
 * 返回：string
 */
func StringAddInt(s string, no int64) string {
	s = strings.TrimSpace(s)

	if no >= 0 && no < 1000 {
		// 加数位数与被加数相等或小于的
		sLen := len(s)
		if sLen < 8 {
			// 整个被加数转成整数
			if is, err := strconv.ParseInt(s, 10, 64); err != nil {
				return "0"
			} else {
				return strconv.FormatInt(is+no, 10)
			}
		} else {
			// 截取比no大一个数量级的数
			save := s[:sLen-8]
			change := s[sLen-8:]
			// 整个被加数转成整数
			if is, err := strconv.ParseInt(change, 10, 64); err != nil {
				return "0"
			} else {
				return save + strconv.FormatInt(is+no, 10)
			}
		}
	} else {
		return s
	}
}
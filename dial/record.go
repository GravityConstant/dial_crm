package dial

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"
)

// 录音
type Record struct {
	Name       string
	PrefixPath string
	File       string
}

func (r *Record) Init(caller, callee string) {
	// 路径前缀
	commonPath := `/home/voices/records`

	now := time.Now()
	year, month, day := now.Date()
	// 年文件夹
	yearStr := strconv.Itoa(year)
	// 月文件夹
	monthStr := strconv.Itoa(int(month))
	// 日文件夹
	dayStr := strconv.Itoa(day)
	// dialplanNumber文件夹
	// 文件名: 主叫-被叫-时间
	nowStr := now.Format("20060102150405000")
	filename := `%s-%s-%s.wav`
	filename = fmt.Sprintf(filename, caller, callee, nowStr)

	filePath := filepath.Join(yearStr, monthStr, dayStr, caller, filename)

	r.Name = "record_file"
	r.PrefixPath = commonPath
	r.File = filePath
}

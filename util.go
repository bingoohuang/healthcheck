package healthcheck

import "time"

// StringSliceContains 检测字符串切片slice是否包含字符串item
func StringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if item == s {
			return true
		}
	}

	return false
}

// MustParseDuration 解析时长字符串，解释失败即panic
func MustParseDuration(s string) time.Duration {
	if duration, err := time.ParseDuration(s); err != nil {
		panic(err)
	} else {
		return duration
	}
}

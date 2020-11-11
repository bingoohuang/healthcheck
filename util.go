package healthcheck

import "time"

// MustParseDuration 解析时长字符串，解释失败即panic
func MustParseDuration(s string) time.Duration {
	if s == "" {
		return 0
	}

	if duration, err := time.ParseDuration(s); err != nil {
		panic(err)
	} else {
		return duration
	}
}

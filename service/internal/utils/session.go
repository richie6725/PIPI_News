package utils

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func GenerateToken() string {
	return uuid.NewString()
}

// CustomTime 是一個自訂的時間型別，用於處理 "YYYY-MM-DD HH:MM:SS" 格式
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		ct.Time = time.Time{}
		return nil
	}

	// 定義可接受的多種時間格式
	layouts := []string{
		time.RFC3339,          // "2025-09-30T23:59:59Z"
		"2006-01-02 15:04:05", // "2025-09-30 23:59:59"
		"2006-01-02",          // "2025-09-30"
	}

	var parsed time.Time
	var err error

	for _, layout := range layouts {
		parsed, err = time.Parse(layout, s)
		if err == nil {
			ct.Time = parsed
			return nil
		}
	}

	// 若全部格式都解析失敗
	return fmt.Errorf("failed to parse time: %s", s)
}

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

const ctLayout = "2006-01-02 15:04:05" // 這是 Go 用來匹配時間格式的特殊記憶字串

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return nil
	}
	t, err := time.Parse(ctLayout, s)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	ct.Time = t
	return nil
}

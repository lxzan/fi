package fi

import "time"

type Valuer interface {
	Value() any
}

// Timestamp 毫秒时间戳
// 会被转化为time.Time类型
type Timestamp int64

func (c Timestamp) Value() any {
	return time.UnixMilli(int64(c))
}

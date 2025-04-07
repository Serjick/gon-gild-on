package gildedtestify

import (
	"time"
)

func init() {
	timeNow = func() time.Time {
		return time.Date(2021, 1, 17, 14, 28, 55, 987654000, time.UTC)
	}
}

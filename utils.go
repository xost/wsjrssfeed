package main

import (
	"time"
)

const timeFormat = "02 Jan 2006 15:04:05"

func transformString2Timestamp(strDate string) time.Time {
	ts, _ := time.Parse(timeFormat, strDate)
	return ts
}

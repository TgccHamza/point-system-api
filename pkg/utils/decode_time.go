package utils

import "time"

// decodeTime decodes a timestamp retrieved from the timeclock
func DecodeTime(t uint32) time.Time {
	second := t % 60
	t = t / 60

	minute := t % 60
	t = t / 60

	hour := t % 24
	t = t / 24

	day := t%31 + 1
	t = t / 31

	month := t%12 + 1
	t = t / 12

	year := t + 2000

	return time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC)
}

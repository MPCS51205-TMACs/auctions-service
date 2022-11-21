package common

import "time"

func InterpretTimeStr(timeStr string) (*time.Time, error) {
	layout := "2006-01-02 15:04:05.000000"
	t, err := time.Parse(layout, timeStr)
	return &t, err
}

func TimeToSQLTimestamp6(aTime time.Time) string {
	layout := "2006-01-02 15:04:05.000000"
	return aTime.Format(layout)
}

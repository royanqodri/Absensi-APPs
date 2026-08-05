package util

import (
	"database/sql"
	"fmt"
	"time"
)

type TimeFormat string
type CustomDate time.Time
type CustomDateTime time.Time

const (
	DATE                      string = "2006-01-02"
	DATETIME                  string = "2006-01-02 15:04:05"
	DATETIME_MS               string = "2006-01-02 15:04:05.000"
	DATETIME_MINUTES_COMBINED string = "200601021504"
	DATE_MINUTES_COMBINED     string = "20060102"
	LOCATION                  string = "Asia/Makassar"
	TIME                      string = "15:04:05"
)

// ToTime converts util.TimeFormat to time.Time
func (tf TimeFormat) ToTime() time.Time {
	return tf.ToTime()
}

func LoadLocation() *time.Location {
	loc, _ := time.LoadLocation(LOCATION)
	return loc
}

func GetTimeNowByLoc() time.Time {
	return time.Now().In(LoadLocation())
}

func GetDateNowByLoc() time.Time {
	t := GetTimeNowByLoc()
	year, month, day := t.Date()
	loc := t.Location()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func GetFormattedDateTimeInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATETIME)
}

func GetFormattedDateTimeMsInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATETIME_MS)
}

func GetFormattedTimeInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TIME)
}

func GetFormattedDateInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATE)
}

func GetFormattedDateTimeMinutesCombinedInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATETIME_MINUTES_COMBINED)
}

func GetFormattedDateMinutesCombinedInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATE_MINUTES_COMBINED)
}

func GetFormattedDate(t string) time.Time {
	parsedTime, err := time.Parse(DATE, t)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func ParseAndFormatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}
	return t.Format("2006-01-02")
}

func GetFormattedDateInLoc(t string) time.Time {
	parsedTime, err := time.ParseInLocation(DATE, t, LoadLocation())
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedDateTime(t string) time.Time {
	parsedTime, err := time.Parse(DATETIME, t)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedDateTimeInLoc(t string) time.Time {
	parsedTime, err := time.ParseInLocation(DATETIME, t, LoadLocation())
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedDateTimeNullable(t *string) time.Time {
	if t == nil || *t == "" {
		return time.Time{}
	}

	parsedTime, err := time.Parse(DATETIME, *t)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedTime(t string) time.Time {
	parsedTime, err := time.Parse(TIME, t)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedTimeWithLayout(t string, layout TimeFormat) time.Time {
	parsedTime, err := time.Parse(string(layout), t)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedTimeWithLayoutInStr(t time.Time, layout TimeFormat) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(string(layout))
}

func GetTimeNowMillisInStr() string {
	now := GetTimeNowByLoc().UnixMilli()
	return IntToStr(now)
}

func (ct *CustomDateTime) UnmarshalText(data []byte) error {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(data))
	if err != nil {
		return err
	}
	*ct = CustomDateTime(parsedTime)
	return nil
}

func (ct CustomDateTime) ToNullTime() sql.NullTime {
	// Assuming that if the time is zero, it's considered invalid
	t := time.Time(ct)
	valid := !t.IsZero()

	return sql.NullTime{
		Time:  t,
		Valid: valid,
	}
}

func (ct *CustomDate) UnmarshalText(data []byte) error {
	parsedTime, err := time.Parse("2006-01-02", string(data))
	if err != nil {
		return err
	}
	*ct = CustomDate(parsedTime)
	return nil
}

func GetDateNowStrID() map[string]string {
	now := GetTimeNowByLoc()

	// Format date
	date := now.Format("2006-01-02 15:04:05.000")
	id := now.Format("20060102150405000")

	// Create map
	objReturn := map[string]string{
		"date": date,
		"id":   id,
	}

	return objReturn
}
func GetFormattedDateStr(t string) time.Time {
	if t == "" {
		return GetTimeNowByLoc()
	}

	parsedTime, err := time.Parse("2006-01-02", t)
	if err != nil {
		return GetTimeNowByLoc()
	}

	return parsedTime
}

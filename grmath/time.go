package grmath

import (
	"encoding/json"

	"time"
)

const (
	DEFAULT_DATETIME = "1970-01-01 00:00:00"
	DATETIME_FORMAT  = "2006-01-02 15:04:05"
)

func GetNanoTimeStampByTimeString(timeString string) (ts int64, err error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return
	}

	tm2, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, loc)
	if err != nil {

		return
	}
	ts = tm2.UnixNano()
	return

}
func GetTimeStampByTimeString(timeString string) (ts int64, err error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return
	}

	tm2, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, loc)
	if err != nil {

		return
	}
	ts = tm2.Unix()
	return

}
func GetCurTimeStamp() (ts int64, err error) {
	nowString := GetFormatCurTime()

	return GetTimeStampByTimeString(nowString)
}

func IsStandardTime(inTime string) (t time.Time, success bool) {
	//timeFormat := "2006-01-02 15:04:05"
	//func Unix(sec int64, nsec int64) Time {

	t, err := time.Parse(DATETIME_FORMAT, inTime)
	if err == nil {
		success = true
	}
	return

}
func TimeToDateTime(t time.Time) string {
	timeStr := t.Format(time.RFC3339)[0:19]
	temp := string(timeStr[0:10])
	temp += " "
	temp += string(timeStr[11:19])
	return temp
}

func GetNowTime() time.Time {
	return time.Now()
}

func GetFormatCurTime() string {

	t := GetNowTime()

	return TimeToDateTime(t)
}

const Fmt = "2006-01-02T15:04:05"

type JTime time.Time

func (jt *JTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {

		return err
	}
	t, err := time.Parse(Fmt, s)
	if err != nil {
		return err
	}
	*jt = (JTime)(t)
	return nil
}
func (jt JTime) MarshalJSON() ([]byte, error) {
	return json.Marshal((*time.Time)(&jt).Format(Fmt))
}

/*
type Time struct {
	time.Time
}

const (
	timeFormart = "2006-01-02 15:04:05"
)

func (t Time) MarshalJSON() ([]byte, error) {

		b := make([]byte, 0, len(timeFormart)+2)
		b = append(b, '"')
		b = time.Time(t).AppendFormat(b, timeFormart)
		b = append(b, '"')


}*/

/*
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}





func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}
*/

package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/jummyliu/pkg/datetime"
)

// Time 实现时间的序列化和反序列化，以及数据库驱动接口
type Time time.Time

// MarshalJSON implements the Marshaler interface
func (t Time) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", datetime.FormatDate(time.Time(t)))
	return []byte(stamp), nil
}

// UnmarshalJSON implements the Unmarshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {
	tmp, err := time.ParseInLocation(fmt.Sprintf("\"%s\"", datetime.DatetimeLayout), string(data), time.Local)
	if err != nil {
		tmp = time.Time{}
	}
	*t = Time(tmp)
	return nil
}

// Scan implements the Scanner interface.
func (t *Time) Scan(src any) error {
	if src == nil {
		return nil
	}
	switch src := src.(type) {
	case string:
		*t = Time(datetime.ParseDate(src))
	case []byte:
		*t = Time(datetime.ParseDate(string(src)))
	case time.Time:
		*t = Time(src)
	default:
		return errors.New("try scan to JSONTime failure")
	}
	return nil
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// String implements flag.Value::String interface.
func (t Time) String() string {
	return datetime.FormatDate(time.Time(t))
}

// Set implements flag.Value::Set interface.
func (t *Time) Set(val string) error {
	tmp, err := time.ParseInLocation(datetime.DatetimeLayout, val, time.Local)
	if err != nil {
		tmp = time.Time{}
	}
	*t = Time(tmp)
	return nil
}

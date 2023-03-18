package envloader

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

type Value interface {
	flag.Value
	flag.Getter
}

func NewString(v string) *String {
	vv := String(v)
	return &vv
}

func StringVar(v *string) *String {
	return (*String)(v)
}

type String string

func (v String) String() string {
	return string(v)
}

func (v String) Get() interface{} {
	return string(v)
}

func (v *String) Set(raw string) error {
	*v = String(raw)
	return nil
}

func NewDuration(v time.Duration) *Duration {
	vv := Duration(v)
	return &vv
}

func DurationVar(v *time.Duration) *Duration {
	return (*Duration)(v)
}

type Duration time.Duration

func (v Duration) String() string {
	return time.Duration(v).String()
}

func (v Duration) Get() interface{} {
	return time.Duration(v)
}

func (v *Duration) Set(raw string) error {
	p, err := time.ParseDuration(raw)
	if err != nil {
		return err
	}
	*v = Duration(p)
	return nil
}

func NewInt(v int) *Int {
	vv := Int(v)
	return &vv
}

func IntVar(v *int) *Int {
	return (*Int)(v)
}

type Int int

func (v Int) String() string {
	return strconv.Itoa(int(v))
}

func (v Int) Get() interface{} {
	return int(v)
}

func (v *Int) Set(raw string) error {
	p, err := strconv.ParseInt(raw, 10, 0)
	if err != nil {
		return err
	}
	*v = Int(p)
	return nil
}

func NewInt64(v int64) *Int64 {
	vv := Int64(v)
	return &vv
}

func Int64Var(v *int64) *Int64 {
	return (*Int64)(v)
}

type Int64 int64

func (v Int64) String() string {
	return strconv.FormatInt(int64(v), 10)
}

func (v Int64) Get() interface{} {
	return int64(v)
}

func (v *Int64) Ptr() *int64 {
	return (*int64)(v)
}

func (v *Int64) Set(raw string) error {
	p, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return err
	}
	*v = Int64(p)
	return nil
}

func NewBool(v bool) *Bool {
	vv := Bool(v)
	return &vv
}

func BoolVar(v *bool) *Bool {
	return (*Bool)(v)
}

type Bool bool

func (v Bool) String() string {
	return strconv.FormatBool(bool(v))
}

func (v Bool) Get() interface{} {
	return bool(v)
}

func (v *Bool) Set(raw string) error {
	p, err := parseBool(raw)
	if err != nil {
		return err
	}
	*v = Bool(p)
	return nil
}

func parseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "on", "On", "ON":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "off", "Off", "OFF":
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean value")
}

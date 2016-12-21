package params

import (
	"fmt"
	"time"
	"bytes"
)

func Aggregate(params []Header) string {
	var buf bytes.Buffer

	for _, param := range params {
		buf.WriteString(param.String())
	}

	return buf.String()
}

type Header interface {
	String() string
}

type stringValue struct {
	key string
	value string
}

func (p stringValue) String() string {
	return fmt.Sprintf(`; %v="%v"`, p.key, p.value)
}

func StringValue(key string, value string) Header {
	return stringValue{
		key: key,
		value: value,
	}
}

type intValue struct {
	key string
	value int
}

func (p intValue) String() string {
	return fmt.Sprintf("; %v=%v", p.key, p.value)
}

func IntValue(key string, value int) Header {
	return intValue{
		key: key,
		value: value,
	}
}

type dateValue struct {
	key string
	value time.Time
}

func (p dateValue) String() string {
	return fmt.Sprintf(`; %v=%v`, p.key, p.value.Format(time.RFC1123))
}

func DateValue(key string, value time.Time) Header {
	return dateValue{
		key: key,
		value: value,
	}
}

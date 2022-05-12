package logger

import (
	"encoding/json"
	"log"
	"time"
)

type Number interface {
	uint8 | int | int64
}

type NewLogger struct {
	WithTimestamp bool
}

func String(key, value string) StringField {
	return StringField{Key: key, Value: value}
}

type StringField struct {
	Key   string
	Value string
}

type FloatField struct {
	Key   string
	Value float64
}

func Float(key string, value float64) FloatField {
	return FloatField{Key: key, Value: value}
}

type IntegerField struct {
	Key   string
	Value int64
}

func Int[V Number](key string, value V) IntegerField {
	val := int64(value)
	return IntegerField{Key: key, Value: val}
}

func (n *NewLogger) Write(message string, fields ...interface{}) (int, error) {
	entryMap := map[string]interface{}{
		"msg": message,
	}

	if n.WithTimestamp {
		entryMap["timestamp"] = time.Now().Format(time.RFC1123Z)
	}

	for _, field := range fields {
		switch f := field.(type) {
		case StringField:
			entryMap[f.Key] = f.Value
		case FloatField:
			entryMap[f.Key] = f.Value
		case IntegerField:
			entryMap[f.Key] = f.Value
		}
	}

	raw, err := json.Marshal(entryMap)
	if err != nil {
		return 0, err
	}

	log.Printf("%s", raw)
	return len(raw), nil
}

package logcore

import "github.com/google/uuid"

type GameLog interface {
	Write(message string, fields ...Field) error
	SetWithTimestamp(t bool)
	SetGameID(id uuid.UUID)
}

type Number interface {
	uint8 | int | int64
}

type Field interface {
	Apply(m map[string]interface{})
}

type StringField struct {
	Key   string
	Value string
}

func String(key, value string) StringField {
	return StringField{Key: key, Value: value}
}

func (s StringField) Apply(m map[string]interface{}) {
	m[s.Key] = s.Value
}

type FloatField struct {
	Key   string
	Value float64
}

func Float(key string, value float64) FloatField {
	return FloatField{Key: key, Value: value}
}

func (s FloatField) Apply(m map[string]interface{}) {
	m[s.Key] = s.Value
}

type IntegerField struct {
	Key   string
	Value int64
}

func Int[V Number](key string, value V) IntegerField {
	val := int64(value)
	return IntegerField{Key: key, Value: val}
}

func (s IntegerField) Apply(m map[string]interface{}) {
	m[s.Key] = s.Value
}

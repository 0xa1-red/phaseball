package logcore

import (
	"bytes"
	"encoding/json"
	"sort"
	"sync"

	"github.com/google/uuid"
)

type GameLog interface {
	Write(message string, fields ...Field) error
	SetWithTimestamp(t bool)
	SetGameID(id uuid.UUID)
	Close() error
}

type Entry struct {
	Seq       int
	Timestamp string
	Entry     map[string]interface{}
}

type EntryCollection struct {
	entries []Entry
	mx      *sync.Mutex
	seq     int
}

func (c *EntryCollection) Add(entry Entry) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.seq += 1
	entry.Seq = c.seq
	entry.Entry["seq"] = c.seq
	b := bytes.NewBuffer([]byte(""))
	encoder := json.NewEncoder(b)
	if err := encoder.Encode(entry.Entry); err != nil {
		return err
	}
	entries := c.entries
	entries = append(entries, entry)
	c.entries = entries
	return nil
}

func NewCollection() *EntryCollection {
	return &EntryCollection{
		entries: make([]Entry, 0),
		mx:      &sync.Mutex{},
		seq:     0,
	}
}

func (c *EntryCollection) Entries() []Entry {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.entries == nil {
		c.entries = make([]Entry, 0)
	}
	entries := c.entries
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Seq > entries[j].Seq
	})
	return entries
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

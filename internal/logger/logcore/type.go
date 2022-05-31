package logcore

import (
	"bytes"
	"encoding/json"
	"sort"
	"sync"

	"github.com/0xa1-red/phaseball/internal/deadball/model"
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

type GameReplay struct {
	ID      uuid.UUID
	Away    model.Team
	Home    model.Team
	Entries *EntryCollection
}

type EntryMap map[string]interface{}

func (e EntryMap) String(key string) string {
	if s, ok := e[key].(string); ok {
		return s
	} else {
		return ""
	}
}

func (e EntryMap) Int64(key string) int64 {
	if s, ok := e[key].(int64); ok {
		return s
	} else {
		return 0
	}
}

type EntryCollection struct {
	entries []Entry
	mx      *sync.Mutex
	seq     int
}

func NewEntryCollection() *EntryCollection {
	return &EntryCollection{
		entries: make([]Entry, 0),
		mx:      &sync.Mutex{},
	}
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
	Apply(m EntryMap)
}

type StringField struct {
	Key   string
	Value string
}

func String(key, value string) StringField {
	return StringField{Key: key, Value: value}
}

func (s StringField) Apply(m EntryMap) {
	m[s.Key] = s.Value
}

type FloatField struct {
	Key   string
	Value float64
}

func Float(key string, value float64) FloatField {
	return FloatField{Key: key, Value: value}
}

func (s FloatField) Apply(m EntryMap) {
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

func (s IntegerField) Apply(m EntryMap) {
	m[s.Key] = s.Value
}

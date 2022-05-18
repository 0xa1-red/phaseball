package psql

import (
	"sync"
	"time"

	"github.com/0xa1-red/phaseball/internal/database"
	"github.com/0xa1-red/phaseball/internal/logger/logcore"
	"github.com/google/uuid"
)

type Logger struct {
	WithTimestamp bool
	GameID        uuid.UUID

	mx      *sync.Mutex
	entries *logcore.EntryCollection
}

func (l *Logger) SetGameID(id uuid.UUID) {
	l.GameID = id
}

func (l *Logger) SetWithTimestamp(t bool) {
	l.WithTimestamp = t
}

func New(opts ...logcore.LoggerOpt) *Logger {
	l := &Logger{
		entries: logcore.NewCollection(),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Logger) Close() error {
	db, err := database.Connection()
	if err != nil {
		return err
	}

	if err := db.WriteGameLog(l.GameID, l.entries.Entries()); err != nil {
		return err
	}
	return nil
}

func (l *Logger) Write(message string, fields ...logcore.Field) error {
	entryMap := map[string]interface{}{
		"msg": message,
	}

	ts := time.Now().Format(time.RFC3339Nano)
	if l.WithTimestamp {
		entryMap["timestamp"] = ts
	}

	for _, field := range fields {
		field.Apply(entryMap)
	}

	return l.entries.Add(logcore.Entry{Timestamp: ts, Entry: entryMap})
}

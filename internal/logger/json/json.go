package json

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/logger/logcore"
	"github.com/google/uuid"
)

type Logger struct {
	WithTimestamp bool
	GameID        uuid.UUID

	seq     int
	entries *logcore.EntryCollection
	fp      *os.File
}

func (l *Logger) SetGameID(id uuid.UUID) {
	l.GameID = id
}

func (l *Logger) SetWithTimestamp(t bool) {
	l.WithTimestamp = t
}

func New(opts ...logcore.LoggerOpt) *Logger {
	path := config.Get().Logging.Path
	if path == "" {
		log.Println("Path is not set for logging, not creating logger facility")
		return nil
	}
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("Error opening path given for logger: %v\n", err)
		return nil
	}

	l := &Logger{
		entries: logcore.NewCollection(),
		fp:      fp,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Logger) Close() error {
	defer l.fp.Close()

	encoder := json.NewEncoder(l.fp)
	return encoder.Encode(l.entries.Entries())
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

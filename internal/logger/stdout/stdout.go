package stdout

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xa1-red/phaseball/internal/logger/logcore"
	"github.com/google/uuid"
)

type Logger struct {
	WithTimestamp bool
	GameID        uuid.UUID
}

func (l *Logger) SetGameID(id uuid.UUID) {
	l.GameID = id
}

func (l *Logger) SetWithTimestamp(t bool) {
	l.WithTimestamp = t
}

func New(opts ...logcore.LoggerOpt) *Logger {
	l := &Logger{}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Logger) Write(message string, fields ...logcore.Field) error {
	entryMap := map[string]interface{}{
		"msg": message,
	}

	ts := time.Now().Format(time.RFC3339)
	if l.WithTimestamp {
		entryMap["timestamp"] = ts
	}

	for _, field := range fields {
		field.Apply(entryMap)
	}

	raw, err := json.Marshal(entryMap)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", raw)
	return nil
}

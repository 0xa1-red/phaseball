package pretty

import (
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

func (l *Logger) Close() error {
	return nil
}

func (l *Logger) Write(message string, fields ...logcore.Field) error {
	var str string
	switch message {
	case logcore.StartOfHalf:
		str = startOfPeriod("half", fields)
	case logcore.StartOfInning:
		str = startOfPeriod("inning", fields)
	case logcore.EndOfHalf:
		str = endOfPeriod("half", fields)
	case logcore.EndOfInning:
		str = endOfPeriod("inning", fields)
	case logcore.AtBat:
		str = atBat(fields)
	default:
		str = message
	}
	if l.WithTimestamp {
		str = fmt.Sprintf("%s - %s", time.Now().Format(time.RFC3339), str)
	}
	fmt.Printf("%s\n", str)
	return nil
}

func startOfPeriod(period string, fields []logcore.Field) string {
	entryMap := map[string]interface{}{
		"msg": fmt.Sprintf("Start of %s", period),
	}

	for _, field := range fields {
		field.Apply(entryMap)
	}

	return fmt.Sprintf("%s: %s %d | Pitching: %s | %d/%d/%d | Pitch die: %s",
		entryMap["msg"].(string),
		entryMap["half"].(string),
		entryMap["inning"].(int64),
		entryMap["pitcher"].(string),
		entryMap["fastball"].(int64),
		entryMap["changeup"].(int64),
		entryMap["breaking"].(int64),
		entryMap["pitch_die"].(string),
	)
}

func endOfPeriod(period string, fields []logcore.Field) string {
	entryMap := map[string]interface{}{
		"msg": fmt.Sprintf("End of %s", period),
	}

	for _, field := range fields {
		field.Apply(entryMap)
	}

	return fmt.Sprintf("%s %d. Hits: %d | Runs: %d\n",
		entryMap["msg"].(string),
		entryMap["inning"].(int64),
		entryMap["hits"].(int64),
		entryMap["runs"].(int64),
	)
}

func atBat(fields []logcore.Field) string {
	entryMap := map[string]interface{}{
		"msg": "At bat",
	}

	for _, field := range fields {
		field.Apply(entryMap)
	}

	return fmt.Sprintf("%s: %s | %d/%d/%d | Batting Target: %d",
		entryMap["msg"].(string),
		entryMap["name"].(string),
		entryMap["power"].(int64),
		entryMap["contact"].(int64),
		entryMap["eye"].(int64),
		entryMap["batter_target"].(int64),
	)
}

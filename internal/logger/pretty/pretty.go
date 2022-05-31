package pretty

import (
	"fmt"
	"time"

	"github.com/0xa1-red/phaseball/internal/deadball/model"
	"github.com/0xa1-red/phaseball/internal/logger/logcore"
	"github.com/google/uuid"
)

var prettyMessages = map[string]string{
	logcore.StartOfHalf:   "Start of half",
	logcore.StartOfInning: "Start of inning",
	logcore.EndOfHalf:     "End of half",
	logcore.EndOfInning:   "End of inning",
	logcore.AtBat:         "At bat",
}

type decoratorFn func(logcore.EntryMap) string

var decorators = map[string]decoratorFn{
	logcore.StartOfHalf:   startOfPeriod,
	logcore.StartOfInning: startOfPeriod,
	logcore.EndOfHalf:     endOfPeriod,
	logcore.EndOfInning:   endOfPeriod,
	logcore.AtBat:         atBat,
	logcore.Pitch:         pitch,
	logcore.Swing:         swing,
	logcore.Out:           out,
	logcore.Hit:           hit,
	logcore.Run:           run,
}

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
	entryMap := logcore.EntryMap{
		"msg": message,
	}

	for _, field := range fields {
		field.Apply(entryMap)
	}

	fn, ok := decorators[message]
	if !ok {
		return fmt.Errorf("Invalid key for decorator: %s", message)
	}

	str := fn(entryMap)

	if l.WithTimestamp {
		str = fmt.Sprintf("%s - %s", time.Now().Format(time.RFC3339), str)
	}
	fmt.Printf("%s\n", str)
	return nil
}

func startOfPeriod(entryMap logcore.EntryMap) string {
	return fmt.Sprintf("%s: %s %d | Pitching: %s | %d/%d/%d | Pitch die: %s",
		prettyMessages[entryMap.String("msg")],
		entryMap.String("half"),
		entryMap.Int64("inning"),
		entryMap.String("pitcher"),
		entryMap.Int64("fastball"),
		entryMap.Int64("changeup"),
		entryMap.Int64("breaking"),
		entryMap.String("pitch_die"),
	)
}

func endOfPeriod(entryMap logcore.EntryMap) string {
	return fmt.Sprintf("%s: %s %d. Hits: %d | Runs: %d\n",
		prettyMessages[entryMap.String("msg")],
		entryMap.String("half"),
		entryMap.Int64("inning"),
		entryMap.Int64("hits"),
		entryMap.Int64("runs"),
	)
}

func atBat(entryMap logcore.EntryMap) string {
	return fmt.Sprintf("%s: %s | %d/%d/%d | Batting Target: %d",
		prettyMessages[entryMap.String("msg")],
		entryMap.String("name"),
		entryMap.Int64("power"),
		entryMap.Int64("contact"),
		entryMap.Int64("eye"),
		entryMap.Int64("batter_target"),
	)
}

func pitch(entryMap logcore.EntryMap) string {
	return fmt.Sprintf("%s is throwing a %s!",
		entryMap.String("pitcher"),
		model.PitchLabels[entryMap["pitch"].(string)],
	)
}

func swing(entryMap logcore.EntryMap) string {
	return fmt.Sprintf("%s swings their bat... (d100: %d; Pitch roll: %d; Pitch modifier: %d; MSS: %d)",
		entryMap.String("name"),
		entryMap.Int64("swing_roll"),
		entryMap.Int64("pitch_roll"),
		entryMap.Int64("pitch_modifier"),
		entryMap.Int64("swing_score"),
	)
}

func out(entryMap logcore.EntryMap) string {
	return fmt.Sprintf("%s is out: %s",
		entryMap.String("name"),
		entryMap.String("event_long"),
	)
}

func hit(entryMap logcore.EntryMap) string {
	event := entryMap.String("event_long")
	if event == "Walk" {
		return fmt.Sprintf("%s walks to first",
			entryMap.String("name"),
		)
	}
	return fmt.Sprintf("%s puts the ball in play: %s",
		entryMap.String("name"),
		event,
	)
}

func run(entryMap logcore.EntryMap) string {
	return fmt.Sprintf("%s comes in for a run!",
		entryMap.String("name"),
	)
}

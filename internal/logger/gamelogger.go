package logger

import (
	"log"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/logger/json"
	"github.com/0xa1-red/phaseball/internal/logger/logcore"
	"github.com/0xa1-red/phaseball/internal/logger/pretty"
	"github.com/0xa1-red/phaseball/internal/logger/psql"
	"github.com/0xa1-red/phaseball/internal/logger/stdout"
	"github.com/google/uuid"
)

const (
	KindStdout string = "stdout"
	KindPsql   string = "psql"
	KindJSON   string = "json"
	KindPretty string = "pretty"
)

func NewGameLogger(id uuid.UUID) logcore.GameLog {
	log.Println(config.Get().GameLog.Kind)
	switch config.Get().GameLog.Kind {
	default:
		fallthrough
	case KindStdout:
		return stdout.New(logcore.WithTimestamp(), logcore.WithGameID(id))
	case KindPsql:
		return psql.New(logcore.WithTimestamp(), logcore.WithGameID(id))
	case KindJSON:
		return json.New(logcore.WithTimestamp(), logcore.WithGameID(id))
	case KindPretty:
		log.Println("Creating 'pretty' logger")
		return pretty.New(logcore.WithTimestamp(), logcore.WithGameID(id))
	}
}

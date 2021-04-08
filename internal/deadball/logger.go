package deadball

import (
	"encoding/json"
	"os"
	"time"

	"github.com/op/go-logging"
	"hq.0xa1.red/axdx/phaseball/internal/logger"
)

var log = getLogger()

func getLogger() *logging.Logger {
	module := "game"

	backend := logging.NewLogBackend(os.Stdout, "", 0)

	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{module} %{longfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	formatter := logging.NewBackendFormatter(backend, format)

	leveled := logging.AddModuleLevel(formatter)

	switch Verbosity() {
	default:
		fallthrough
	case verboseWarning:
		leveled.SetLevel(logging.WARNING, module)
	case verboseInfo:
		leveled.SetLevel(logging.INFO, module)
	case verboseDebug:
		leveled.SetLevel(logging.DEBUG, module)

	}

	log := logger.New(module)
	log.SetBackend(leveled)

	return log
}

type GameLog struct {
	PlayedAt time.Time
	Away     *Team
	Home     *Team
	Entries  []LogEntry
}

func NewGameLog(away, home *Team) *GameLog {
	return &GameLog{
		PlayedAt: time.Now(),
		Away:     away,
		Home:     home,
		Entries:  make([]LogEntry, 0),
	}
}

func (g *GameLog) Append(l LogEntry) {
	g.Entries = append(g.Entries, l)
}

func (g *GameLog) String() string {
	if j, err := json.MarshalIndent(g, "", "    "); err != nil {
		return ""
	} else {
		return string(j)
	}
}

type LogEntry struct {
	Inning  *InningLog
	Batter  *Player
	Pitcher *Player
	Event   Event
	Runs    int
	Extra   map[string]interface{}
	Bases   struct {
		First  string
		Second string
		Third  string
		Home   string
	}
}

func (e LogEntry) JSON() (string, error) {
	j, err := json.Marshal(e)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

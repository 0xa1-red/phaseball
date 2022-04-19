package deadball

import (
	"encoding/json"
	"os"
	"time"

	"github.com/0xa1-red/phaseball/internal/logger"
	"github.com/op/go-logging"
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

type Score struct {
	Innings map[int]int
	Hits    uint8
	Runs    uint8
}

type GameLog struct {
	PlayedAt time.Time
	Away     *Team
	Home     *Team
	Entries  []LogEntry
	BoxScore map[string]Score
}

func NewGameLog(away, home *Team) *GameLog {
	return &GameLog{
		PlayedAt: time.Now(),
		Away:     away,
		Home:     home,
		Entries:  make([]LogEntry, 0),
		BoxScore: make(map[string]Score),
	}
}

func (g *GameLog) Append(l LogEntry) {
	g.Entries = append(g.Entries, l)
}

func (g *GameLog) String() string {
	if j, err := json.MarshalIndent(g, "", "    "); err != nil {
		panic(err)
	} else {
		return string(j)
	}
}

func (g *GameLog) AddInning(inning int, team string, hits, runs uint8) {
	if _, ok := g.BoxScore[team]; !ok {
		g.BoxScore[team] = Score{
			Innings: make(map[int]int, 0),
			Hits:    0,
			Runs:    0,
		}
	}

	score := g.BoxScore[team]
	score.Innings[inning] = int(runs)
	score.Hits += hits
	score.Runs += runs
	g.BoxScore[team] = score
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

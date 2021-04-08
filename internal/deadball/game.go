package deadball

import (
	"fmt"

	"hq.0xa1.red/axdx/phaseball/internal/dice"
)

/**
Game represents a single baseball match.
A game runs for at least 9 turns consisting of a top and a bottom inning.
If the score (runs) is even after 9 turns, the game goes into overtime, meaning the first full
turn with a winning score ends the game.

TODO: The game should end if the team playing the bottom inning in the last round is in the lead
*/
type Game struct {
	Turns []*Turn
	Teams map[string]*Team
	Log   *GameLog
}

func (g *Game) NewInning(num uint8, half string) *Inning {
	switch half {
	case HalfTop:
		return NewInning(g.Teams[TeamAway], g.Teams[TeamHome].Pitcher(), g.Log, num, half)
	case HalfBottom:
		return NewInning(g.Teams[TeamHome], g.Teams[TeamAway].Pitcher(), g.Log, num, half)
	}

	return &Inning{}
}

func NewGame(away, home Team) *Game {
	return &Game{
		Turns: make([]*Turn, 0),
		Teams: map[string]*Team{
			TeamAway: &away,
			TeamHome: &home,
		},
	}
}

func (g *Game) Score() map[string]uint8 {
	var home, away uint8
	for _, turn := range g.Turns {
		away += turn.Top.Runs
		home += turn.Bottom.Runs
	}
	return map[string]uint8{
		TeamAway: away,
		TeamHome: home,
	}
}

// Run simulates the game from start to finish
func (g *Game) Run() {
	g.Log = NewGameLog(g.Teams[TeamAway], g.Teams[TeamHome])
	for i := 0; i < 9; i++ {
		log.Debugf("Inning %d\n", i+1)
		turn := &Turn{
			Top:    g.NewInning(uint8(i+1), HalfTop),
			Bottom: g.NewInning(uint8(i+1), HalfBottom),
		}

		log.Debugf("Inning %d - TOP - %s\n", i+1, g.Teams[TeamAway].Name)
		turn.Top.Run()

		log.Debugf("Inning %d - BOTTOM - %s\n", i+1, g.Teams[TeamHome].Name)
		turn.Bottom.Run()
		g.Turns = append(g.Turns, turn)
	}

	runs := g.Score()
	awayRuns := runs[TeamAway]
	homeRuns := runs[TeamHome]

	if awayRuns == homeRuns {
		i := 9
		for awayRuns == homeRuns {
			log.Debugf("Inning %d\n", i+1)
			turn := &Turn{
				Top:    g.NewInning(uint8(i+1), HalfTop),
				Bottom: g.NewInning(uint8(i+1), HalfBottom),
			}

			log.Debugf("Inning %d - TOP - %s\n", i+1, g.Teams[TeamAway].Name)
			turn.Top.Run()

			log.Debugf("Inning %d - BOTTOM - %s\n", i+1, g.Teams[TeamHome].Name)
			turn.Bottom.Run()
			g.Turns = append(g.Turns, turn)

			runs := g.Score()
			awayRuns = runs[TeamAway]
			homeRuns = runs[TeamHome]
			i++
		}
	}
}

// Turn in lack of a better term represents a top and a bottom inning
type Turn struct {
	Top    *Inning
	Bottom *Inning
}

// Inning represents a teams turn at batting
type Inning struct {
	Number  uint8
	Outs    uint8
	Hits    uint8
	Runs    uint8
	Half    string
	Logger  *GameLog
	Team    *Team
	Pitcher *Player
	Diamond *Diamond
}

// NewInning creates a new inning for a team
func NewInning(team *Team, pitcher *Player, log *GameLog, num uint8, half string) *Inning {
	inning := Inning{
		Team:    team,
		Pitcher: pitcher,
		Number:  num,
		Half:    half,
		Logger:  log,
		Diamond: GetDiamond(),
	}

	inning.Team.NewTurn(false)
	if num > 1 {
		inning.Team.Next()
	}
	inning.Diamond.Reset()

	return &inning
}

// Run simulates an inning
func (i *Inning) Run() {
	i.Team.NewTurn(false)
	for i.Outs < 3 {
		i.AtBat()
	}
}

func (i *Inning) ToLog() *InningLog {
	return &InningLog{
		Number: i.Number,
		Half:   i.Half,
		Outs:   i.Outs,
		Runs:   i.Runs,
		Hits:   i.Hits,
	}
}

type InningLog struct {
	Number uint8
	Half   string
	Hits   uint8
	Outs   uint8
	Runs   uint8
}

func (i *Inning) ProductiveOut(swing int, outEvent Event) (Event, []*Player) {
	runners := make([]*Player, 0)
	event := outEvent
	if swing < 70 && IsOutOutfield(outEvent) && i.Outs < 3 {
		if p2 := i.Diamond.Bases[1].Player; p2 != nil {
			if pp := i.Diamond.Bases[1].Load(nil); pp != nil {
				runners = append(runners, pp)
			}
			event = Event{Label: EventHitProductiveOut.Label, Extra: event.Extra}
		} else if i.Diamond.Bases[2].Player != nil {
			if runner := i.Diamond.Bases[2].Load(nil); runner != nil {
				runners = append(runners, runner)
			}
			event = Event{Label: EventHitProductiveOut.Label, Extra: event.Extra}
		}
	}

	return event, runners
}

func (i *Inning) PossibleDouble(swing int, outEvent Event, p *Player) Event {
	event := outEvent
	digit := LastDigit(swing)
	if i.Diamond.Bases[0].Player != nil && IsOutInfield(outEvent) && digit >= 3 && digit < 7 {
		if swing >= 70 {
			i.Diamond.Bases[0].Player = nil
			event = Event{Label: EventHitDoublePlay.Label, Extra: event.Extra}
			i.Outs++
		} else {
			i.Diamond.Bases[0].Player = p
			event = Event{Label: EventHitFieldersChoice.Label, Extra: event.Extra}
		}
	}

	return event
}

// AtBat simulates a single player's at bat scenario
func (i *Inning) AtBat() {
	// Select the batter
	p := i.Team.AtBat()

	l := LogEntry{
		Inning:  i.ToLog(),
		Batter:  p,
		Pitcher: i.Pitcher,
		Extra:   make(map[string]interface{}),
	}

	_, pitchRoll := i.Pitcher.Pitch(p.Hand)

	swing := dice.Roll(100, 1, 0) + pitchRoll
	event := swingEvent(swing, int(p.BatterTarget))

	var scored int
	runners := []*Player{}
	switch event {
	case EventProdOut, EventPossibleDbl:
		p.Status = StatusOut
		i.Outs++
	case EventHit, EventCrit:
		hitResult, extra, out := Hit(swing, event == EventCrit)

		if extra {
			if runner := i.Diamond.Bases[0].Load(nil); runner != nil {
				runners = append(runners, runner)
				scored++
			}
		}

		l.Event = hitResult

		switch hitResult {
		case EventHitSinglePlus:
			if runs := i.Diamond.Single(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
		case EventHitSingleAdv2:
			if runner := i.Diamond.Bases[0].Load(nil); runner != nil {
				runners = append(runners, runner)
				scored++
			}
			if runs := i.Diamond.Single(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
		case EventHitDoubleAdv3:
			if i.Diamond.Bases[0].Load(nil) != nil {
				scored++
			}
			if i.Diamond.Bases[0].Load(nil) != nil {
				scored++
			}
			if runs := i.Diamond.Double(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
		case EventHitHomeRun:
			if runs := i.Diamond.HomeRun(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			i.Team.Next()
		case EventHitSingle1B, EventHitSingle2B, EventHitSingle3B, EventHitSingleSS:
			if out {
				p.Status = StatusOut
				i.Outs++
				break
			}
			if runs := i.Diamond.Single(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			if extra {
				l.Event = EventHitSingleError
			} else {
				l.Event = EventHitSingle
			}
		case EventHitDoubleCF, EventHitDoubleLF, EventHitDoubleRF:
			if out {
				l.Event = hitResult
				p.Status = StatusOut
				i.Outs++
				break
			}
			if runs := i.Diamond.Double(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			if extra {
				l.Event = EventHitDoubleError
			} else {
				l.Event = EventHitDouble
			}
		case EventHitTripleCF, EventHitTripleRF:
			if out {
				l.Event = hitResult
				p.Status = StatusOut
				i.Outs++
				break
			}
			if runs := i.Diamond.Triple(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			if extra {
				l.Event = EventHitTripleError
			} else {
				l.Event = EventHitTriple
			}
		}

		i.Hits++
	case EventWalk:
		l.Event = EventWalk // TODO: Refactor events and rename this
		if runs := i.Diamond.Advance(p, 1); len(runs) > 0 {
			runners = append(runners, runs...)
			scored += len(runs)
		}
	}

	if p.Status == StatusOut {
		outEvent := Out(swing)
		l.Event = outEvent
		if event == EventProdOut {
			if event, prunners := i.ProductiveOut(swing, outEvent); event == EventHitProductiveOut {
				runners = append(runners, prunners...)
				l.Event = event
			}
		} else if event == EventPossibleDbl && i.Outs < 2 {
			outEvent := i.PossibleDouble(swing, outEvent, p)
			l.Event = outEvent
		}
	}

	if Verbosity() == verboseDebug {
		log.Debug(i.Diamond.String())
	}

	if len(runners) > 0 {
		l.Extra["runners"] = runners
		l.Runs = int(scored)
		i.Runs = i.Runs + uint8(scored)
		log.Debugf("\t\tRBI: %d\n", scored)
	}

	bases := struct {
		First  string
		Second string
		Third  string
		Home   string
	}{}
	for _, base := range i.Diamond.Bases {
		if base.Player != nil {
			switch base.Name {
			case BaseFirst:
				bases.First = base.Player.Name
			case BaseSecond:
				bases.Second = base.Player.Name
			case BaseThird:
				bases.Third = base.Player.Name
			case BaseHome:
				bases.Home = base.Player.Name
			}
		}
	}
	l.Bases = bases

	i.Logger.Append(l)
}

// diamond is a singleton holding the current state of the pitch
var diamond *Diamond

// Diamond holds the three bases
type Diamond struct {
	Bases [4]*Base
}

func (d *Diamond) String() string {
	str := ""
	for _, base := range d.Bases {
		name := "empty"
		if base.Player != nil {
			name = base.Player.Name
		}
		str = fmt.Sprintf("%s\t%s : %s\n", str, base.Name, name)
	}

	return str
}

// Reset creates a new diamond with empty bases
func (d *Diamond) Reset() {
	diamond = newDiamond()
}

// newDiamond creates a new diamond with empty bases
func newDiamond() *Diamond {
	var bases [4]*Base
	baseStr := []string{
		BaseFirst,
		BaseSecond,
		BaseThird,
		BaseHome,
	}
	for i := 0; i < 4; i++ {
		bases[i] = &Base{
			Name: baseStr[i],
		}
		if i > 0 {
			bases[i-1].Next = bases[i]
		}
	}
	return &Diamond{
		Bases: bases,
	}
}

// GetDiamond returns the current state or creates a new one if it's empty
func GetDiamond() *Diamond {
	if diamond == nil {
		diamond = newDiamond()
	}

	return diamond
}

// Advance pushes a batter up n bases
func (d *Diamond) Advance(p *Player, bases int) []*Player {
	runs := make([]*Player, 0)
	if bases == 4 {
		for i := len(d.Bases) - 1; i >= 0; i-- {
			// For every player on base we add a run, reset their status and remove from the diamond
			if d.Bases[i].Player != nil {
				runs = append(runs, d.Bases[i].Player)
				d.Bases[i].Player.Status = StatusWaiting
				d.Bases[i].Player = nil
			}
		}
		// We add a run for the batting player
		runs = append(runs, p)

		return runs
	}
	p.Status = StatusBase
	for i := 0; i < bases; i++ {
		if i > 0 {
			p = nil
		}
		if runner := d.Bases[0].Load(p); runner != nil {
			runs = append(runs, runner)
		}
	}

	return runs
}

// Single is a shorthand for advancing a base
func (d *Diamond) Single(p *Player) []*Player {
	return d.Advance(p, 1)
}

// Double is a shorthand for advancing 2 bases
func (d *Diamond) Double(p *Player) []*Player {
	return d.Advance(p, 2)
}

// Triple is a shorthand for advancing 3 bases
func (d *Diamond) Triple(p *Player) []*Player {
	return d.Advance(p, 3)
}

// HomeRun is a shorthand for running in and clearing all loaded bases
func (d *Diamond) HomeRun(p *Player) []*Player {
	return d.Advance(p, 4)
}

// Base represents a single base linked to the next and holding whether a player is standing on it
type Base struct {
	Name   string
	Next   *Base
	Player *Player
}

// Load puts a player on a base or in some cases pushes a player further from that base
// without loading it.
func (b *Base) Load(p *Player) *Player {
	pp := "empty"
	if p != nil {
		pp = p.Name
	}
	log.Debugf("Loading %s to %s", pp, b.Name)

	if b.Name == BaseHome {
		log.Debugf("home base")
		return p
	}

	if b.Player == nil && p == nil {
		log.Debugf("player nil, next nil")
		return b.Next.Load(nil)
	}

	if b.Player == nil && p != nil {
		log.Debugf("player nil, next not nil")
		b.Player = p
		return b.Next.Load(nil)
	}

	if b.Player != nil && p == nil {
		log.Debugf("player not nil, next nil")
		run := b.Next.Load(b.Player)
		b.Player = nil
		return run
	}

	if b.Player != nil && p != nil {
		log.Debugf("player not nil, next not nil")
		run := b.Next.Load(b.Player)
		b.Player = p
		return run
	}

	return nil
}

func swingEvent(swing int, bt int) Event {
	var event Event
	if swing >= 71 {
		event = EventPossibleDbl
	} else if swing >= bt+6 {
		event = EventProdOut
	} else if swing >= bt+1 {
		event = EventWalk
	} else if swing >= 6 {
		event = EventHit
	} else {
		event = EventCrit
	}

	return event
}

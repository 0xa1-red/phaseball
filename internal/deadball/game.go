package deadball

import (
	"fmt"
	"strings"

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
	Log   GameLog
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
	for i := 0; i < 9; i++ {
		log.Debugf("Inning %d\n", i+1)
		turn := &Turn{
			Top:    NewInning(g.Teams[TeamAway], g.Teams[TeamHome].Pitcher(), uint8(i+1), HalfTop),
			Bottom: NewInning(g.Teams[TeamHome], g.Teams[TeamAway].Pitcher(), uint8(i+1), HalfBottom),
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
				Top:    NewInning(g.Teams[TeamAway], g.Teams[TeamHome].Pitcher(), uint8(i+1), HalfTop),
				Bottom: NewInning(g.Teams[TeamHome], g.Teams[TeamAway].Pitcher(), uint8(i+1), HalfBottom),
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

	header := fmt.Sprintf("%-30s |", "Team name")
	separator := fmt.Sprintf("%s+", strings.Repeat("-", 31))
	awayBoard := fmt.Sprintf("%-30s |", g.Teams[TeamAway].Name)
	homeBoard := fmt.Sprintf("%-30s |", g.Teams[TeamHome].Name)

	var awayHits uint8 = 0
	var homeHits uint8 = 0

	for i, turn := range g.Turns {
		awayHits += turn.Top.Hits
		homeHits += turn.Bottom.Hits

		header = fmt.Sprintf("%s %4d |", header, i+1)
		separator = fmt.Sprintf("%s%s+", separator, strings.Repeat("-", 6))
		awayBoard = fmt.Sprintf("%s %4d |", awayBoard, turn.Top.Runs)
		homeBoard = fmt.Sprintf("%s %4d |", homeBoard, turn.Bottom.Runs)
	}

	header = fmt.Sprintf("%s| %4s | ", header, "H")
	header = fmt.Sprintf("%s %4s | ", header, "R")
	separator = fmt.Sprintf("%s+%s+%s+", separator, strings.Repeat("-", 6), strings.Repeat("-", 7))

	awayBoard = fmt.Sprintf("%s| %4d | ", awayBoard, awayHits)
	awayBoard = fmt.Sprintf("%s %4d | ", awayBoard, awayRuns)

	homeBoard = fmt.Sprintf("%s| %4d | ", homeBoard, homeHits)
	homeBoard = fmt.Sprintf("%s %4d | ", homeBoard, homeRuns)

	fmt.Println(header)
	fmt.Println(separator)
	fmt.Println(awayBoard)
	fmt.Println(homeBoard)

	g.Log = gameLog
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
	Team    *Team
	Pitcher *Player
	Diamond *Diamond
}

// NewInning creates a new inning for a team
func NewInning(team *Team, pitcher *Player, num uint8, half string) *Inning {
	inning := Inning{
		Team:    team,
		Pitcher: pitcher,
		Number:  num,
		Half:    half,
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
	}
}

type InningLog struct {
	Number uint8
	Half   string
	Outs   uint8
	Runs   uint8
}

func (i *Inning) ProductiveOut(swing int, outEvent Event) (ExtendedEvent, []*Player) {
	runners := make([]*Player, 0)
	event := ExtendedEventMapping[outEvent]
	if swing < 70 && IsOutOutfield(outEvent) && i.Outs < 3 {
		if p2 := i.Diamond.Bases[1].Player; p2 != nil {
			if pp := i.Diamond.Bases[1].Load(nil); pp != nil {
				runners = append(runners, pp)
			}
			event = ExtendedEvent{EventHitProductiveOut, event.Extra}
		} else if i.Diamond.Bases[2].Player != nil {
			if runner := i.Diamond.Bases[2].Load(nil); runner != nil {
				runners = append(runners, runner)
			}
			event = ExtendedEvent{EventHitProductiveOut, event.Extra}
		}
	}

	return event, runners
}

func (i *Inning) PossibleDouble(swing int, outEvent Event, p *Player) ExtendedEvent {
	event := ExtendedEventMapping[outEvent]
	digit := LastDigit(swing)
	if i.Diamond.Bases[0].Player != nil && IsOutInfield(outEvent) && digit >= 3 && digit < 7 {
		if swing >= 70 {
			i.Diamond.Bases[0].Player = nil
			event = ExtendedEvent{EventHitDoublePlay, event.Extra}
			i.Outs++
		} else {
			i.Diamond.Bases[0].Player = p
			event = ExtendedEvent{EventHitFieldersChoice, event.Extra}
		}
	}

	return event
}

// AtBat simulates a single player's at bat scenario
func (i *Inning) AtBat() {
	// Select the batter
	p := i.Team.AtBat()
	events := map[int]string{
		EventCrit:        EventCritStr,
		EventHit:         EventHitStr,
		EventWalk:        EventWalkStr,
		EventProdOut:     EventProdOutStr,
		EventPossibleDbl: EventPossibleDblStr,
	}

	log.Debugf("Pitcher: %s | Batter: %s", i.Pitcher.Name, p.Name)

	l := LogEntry{
		Inning:  i.ToLog(),
		Batter:  p,
		Pitcher: i.Pitcher,
		Extra:   make(map[string]interface{}),
	}

	var eventKey int

	_, pitchRoll := i.Pitcher.Pitch(p.Hand)

	swing := dice.Roll(100, 1, 0) + pitchRoll
	str := fmt.Sprintf("BT: %d | SS: %d", p.BatterTarget, swing)
	if swing >= 71 {
		str = fmt.Sprintf("%s | Result: %s", str, "Possible double")
		eventKey = EventPossibleDbl
	} else if swing >= int(p.BatterTarget)+6 {
		str = fmt.Sprintf("%s | Result: %s", str, "Productive out")
		eventKey = EventProdOut
	} else if swing >= int(p.BatterTarget)+1 {
		str = fmt.Sprintf("%s | Result: %s", str, "Walk")
		eventKey = EventWalk
	} else if swing >= 6 {
		str = fmt.Sprintf("%s | Result: %s", str, "Hit")
		eventKey = EventHit
	} else {
		str = fmt.Sprintf("%s | Result: %s", str, "Critical")
		eventKey = EventCrit
	}
	log.Debug(str)

	eventStr := events[eventKey]

	var scored int
	runners := []*Player{}
	switch eventKey {
	case EventProdOut, EventPossibleDbl:
		p.Status = StatusOut
		i.Outs++
	case EventHit, EventCrit:
		hitResult, extra, out := Hit(swing, eventKey == EventCrit)

		if extra {
			if runner := i.Diamond.Bases[0].Load(nil); runner != nil {
				runners = append(runners, runner)
				scored++
			}
		}

		l.Event = ExtendedEvent{hitResult, ""}

		switch hitResult {
		case EventHitSinglePlus:
			if runs := i.Diamond.Single(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			eventStr = fmt.Sprintf("%s - Single", eventStr)
		case EventHitSingleAdv2:
			if runner := i.Diamond.Bases[0].Load(nil); runner != nil {
				runners = append(runners, runner)
				scored++
			}
			if runs := i.Diamond.Single(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			eventStr = fmt.Sprintf("%s - Single, runners advence 2", eventStr)
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
			eventStr = fmt.Sprintf("%s - Double, runners advance 3", eventStr)
		case EventHitHomeRun:
			if runs := i.Diamond.HomeRun(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			eventStr = fmt.Sprintf("%s - Home run!!!", eventStr)
		case EventHitSingle1B, EventHitSingle2B, EventHitSingle3B, EventHitSingleSS:
			positions := map[Event]string{
				EventHitSingle1B: "first baseman",
				EventHitSingle2B: "second baseman",
				EventHitSingle3B: "third baseman",
				EventHitSingleSS: "shortstop",
			}
			eventStr = fmt.Sprintf("Defender: %s", positions[hitResult])
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
				l.Event = ExtendedEvent{EventHitSingleError, ""}
				eventStr = fmt.Sprintf("%s | Single, Error", eventStr)
			} else {
				l.Event = ExtendedEvent{EventHitSingle, ""}
				eventStr = fmt.Sprintf("%s | Single", eventStr)
			}
		case EventHitDoubleCF, EventHitDoubleLF, EventHitDoubleRF:
			positions := map[Event]string{
				EventHitDoubleCF: "center fielder",
				EventHitDoubleLF: "left fielder",
				EventHitDoubleRF: "right fielder",
			}
			eventStr = fmt.Sprintf("Defender: %s", positions[hitResult])
			if out {
				l.Event = ExtendedEvent{hitResult, ""}
				p.Status = StatusOut
				i.Outs++
				break
			}
			if runs := i.Diamond.Double(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			if extra {
				l.Event = ExtendedEvent{EventHitDoubleError, ""}
				eventStr = fmt.Sprintf("%s | Double, Error", eventStr)
			} else {
				l.Event = ExtendedEvent{EventHitDouble, ""}
				eventStr = fmt.Sprintf("%s | Double", eventStr)
			}
		case EventHitTripleCF, EventHitTripleRF:
			positions := map[Event]string{
				EventHitTripleCF: "center fielder",
				EventHitTripleRF: "right fielder",
			}
			eventStr = fmt.Sprintf("Defender: %s", positions[hitResult])
			if out {
				l.Event = ExtendedEvent{hitResult, ""}
				p.Status = StatusOut
				i.Outs++
				break
			}
			if runs := i.Diamond.Triple(p); len(runs) > 0 {
				runners = append(runners, runs...)
				scored += len(runs)
			}
			if extra {
				l.Event = ExtendedEvent{EventHitTripleError, ""}
				eventStr = fmt.Sprintf("%s | Triple, Error", eventStr)
			} else {
				l.Event = ExtendedEvent{EventHitTriple, ""}
				eventStr = fmt.Sprintf("%s | Triple", eventStr)
			}
		}

		i.Hits++
	case EventWalk:
		if runs := i.Diamond.Advance(p, 1); len(runs) > 0 {
			l.Event = ExtendedEvent{EventLogWalk, ""} // TODO: Refactor events and rename this
			runners = append(runners, runs...)
			scored += len(runs)
		}
	}

	if p.Status == StatusOut {
		outEvent := Out(swing)
		l.Event = ExtendedEventMapping[outEvent]
		if eventKey == EventProdOut {
			if event, prunners := i.ProductiveOut(swing, outEvent); event.Event == EventHitProductiveOut {
				runners = append(runners, prunners...)
				outEvent = event.Event
				l.Event = event
			}
			eventStr = fmt.Sprintf("%s | %s", eventStr, outEvent.Long())
		} else if eventKey == EventPossibleDbl {
			outEvent := i.PossibleDouble(swing, outEvent, p)
			l.Event = outEvent

			eventStr = fmt.Sprintf("%s | %s", eventStr, outEvent.Event.Long())
		} else {
			eventStr = fmt.Sprintf("%s | %s", eventStr, outEvent.Long())
		}
	}

	eventOutput := fmt.Sprintf("At bat: %s ... %s", p.Name, eventStr)
	// Send information and debug messages to the logger
	log.Debug(eventOutput)

	for _, base := range i.Diamond.Bases {
		name := "empty"
		if base.Player != nil {
			name = base.Player.Name
		}
		log.Debugf("\t%s : %s\n", base.Name, name)
	}

	if scored > 0 {
		l.Extra["runners"] = runners
		l.Runs = int(scored)
		i.Runs = i.Runs + uint8(scored)
		log.Debugf("\t\tRBI: %d\n", scored)
	}

	gameLog = append(gameLog, l)
}

// diamond is a singleton holding the current state of the pitch
var diamond *Diamond

// Diamond holds the three bases
type Diamond struct {
	Bases [4]*Base
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

package deadball

import (
	"fmt"
	"strings"

	"github.com/0xa1-red/phaseball/internal/dice"
)

const (
	HandLeftie  = "LEFT"
	HandRightie = "RIGHT"
	HandSwitch  = "SWITCH"
)

var (
	Pitcher   = Position{ID: PositionPitcher, Short: "P", Name: "Pitcher"}
	Catcher   = Position{ID: PositionCatcher, Short: "C", Name: "Catcher"}
	First     = Position{ID: PositionFirst, Short: "1B", Name: "First baseman"}
	Second    = Position{ID: PositionSecond, Short: "2B", Name: "Second baseman"}
	Third     = Position{ID: PositionThird, Short: "3B", Name: "Third baseman"}
	Shortstop = Position{ID: PositionShortstop, Short: "SS", Name: "Shortstop"}
	Left      = Position{ID: PositionLeft, Short: "LF", Name: "Left fielder"}
	Center    = Position{ID: PositionCenter, Short: "CF", Name: "Center fielder"}
	Right     = Position{ID: PositionRight, Short: "RF", Name: "Right fielder"}
)

// Team represents a team in a match
type Team struct {
	Name    string
	Players [9]*Player
	Index   int `json:"-"`
}

// NewTurn resets the players' status in a team
func (t *Team) NewTurn(reindex bool) {
	for _, player := range t.Players {
		player.Status = StatusWaiting
	}
	if reindex {
		t.Index = 0
	}
}

// AtBat returns the player who's supposed to bat next
func (t *Team) AtBat() *Player {
	var next *Player
	for next == nil {
		if t.Players[t.Index].Status == StatusWaiting {
			next = t.Players[t.Index]
			break
		}
		t.Next()
	}

	return next
}

// OnDeck returns the player who's supposed to bat after the current batter
func (t *Team) OnDeck() *Player {
	i := t.Index
	var next *Player
	for next == nil {
		if t.Players[i].Status == StatusWaiting {
			next = t.Players[i]
			break
		}
		i++
		if i > 8 {
			i = 0
		}
	}

	return next
}

// Next increases the index until it reaches the maximum then resets to 0
func (t *Team) Next() {
	if i := t.Index + 1; i > 8 {
		t.Index = 0
	} else {
		t.Index = i
	}
}

func (t *Team) Pitcher() *Player {
	for _, player := range t.Players {
		if player.Position.ID == PositionPitcher {
			return player
		}
	}
	return nil
}

// String returns the string representation of the team
func (t *Team) String() string {
	strs := make([]string, 0)
	for i, player := range t.Players {
		status := player.Status
		if i == t.Index {
			status = StatusOnDeck
		}

		strs = append(strs, fmt.Sprintf("\t%s (%s)", player.Name, status))
	}

	return strings.Join(strs, "\n")
}

type Position struct {
	ID    byte
	Short string
	Name  string
}

// Player represents a single player
type Player struct {
	Power        int
	Contact      int
	Eye          int
	Speed        int
	Defense      int
	Fastball     int
	Changeup     int
	Breaking     int
	Control      int
	Batting      int
	Name         string
	Status       string `json:"-"`
	Position     Position
	BatterTarget uint8
	PitchDie     PitchDie
	Hand         string
	Traits       []Trait
}

// NewPlayer returns a new player
func NewPlayer(name string, position Position) *Player {
	return &Player{
		Name:     name,
		Status:   StatusWaiting,
		Position: position,
	}
}

func (p *Player) Pitch(batterHand string) (PitchDie, int) {
	pitcherAdvantage := p.Hand == batterHand

	if p.PitchDie == PitchAddD12 || pitcherAdvantage && p.PitchDie == PitchAddD8 {
		return PitchAddD12, dice.Roll(12, 1, 0)
	} else if p.PitchDie == PitchAddD8 || pitcherAdvantage && p.PitchDie == PitchAddD4 {
		return PitchAddD8, dice.Roll(8, 1, 0)
	} else if p.PitchDie == PitchAddD4 || pitcherAdvantage && p.PitchDie == PitchSubD4 {
		return PitchAddD4, dice.Roll(4, 1, 0)
	} else if p.PitchDie == PitchSubD4 {
		return PitchSubD4, dice.Roll(4, 1, 0) * -1
	}

	return PitchNone, 0
}

func (p *Player) CalculateDie() {
	avg := float64(p.Fastball+p.Changeup+p.Breaking+p.Control) / 4

	if avg <= 5 {
		p.PitchDie = PitchSubD4
	} else if avg <= 6 {
		p.PitchDie = PitchAddD4
	} else if avg <= 7 {
		p.PitchDie = PitchAddD8
	} else {
		p.PitchDie = PitchAddD12
	}
}

func (p *Player) HasTrait(t Trait) bool {
	for _, trait := range p.Traits {
		if t == trait {
			return true
		}
	}

	return false
}

func (p *Player) Hit(swing int, crit bool) (result Event, extra bool, out bool) {
	roll := dice.Roll(20, 1, 0)

	if p.HasTrait(TraitPower) {
		roll += 1
	} else if p.HasTrait(TraitPowerPlus) {
		roll += 2
	} else if p.HasTrait(TraitWeak) {
		roll -= 1
	} else if p.HasTrait(TraitWeakMinus) {
		roll -= 2
	}

	return p.hit(swing, crit, roll)
}

func (p *Player) hit(swing int, crit bool, roll int) (Event, bool, bool) {
	var (
		result Event
		extra  bool
		out    bool
	)

	if roll >= 19 || crit && roll == 18 {
		result, extra, out = EventHitHomeRun, false, false
	} else if roll == 18 || crit && roll >= 16 {
		if swing%2 == 0 {
			result, extra, out = Defense(EventHitTripleRF)
		} else {
			result, extra, out = Defense(EventHitTripleCF)
		}
	} else if roll >= 16 || crit && roll == 15 {
		result, extra, out = EventHitDoubleAdv3, false, false
	} else if roll == 15 || crit && roll == 14 {
		result, extra, out = Defense(EventHitDoubleRF)
	} else if roll == 14 || crit && roll == 13 {
		result, extra, out = Defense(EventHitDoubleCF)
	} else if roll == 13 || crit && roll >= 8 {
		result, extra, out = Defense(EventHitDoubleLF)
	} else if roll >= 8 || crit && roll == 7 {
		result, extra, out = EventHitSingleAdv2, false, false
	} else if roll == 7 || crit && roll == 6 {
		if swing%2 == 0 {
			result, extra, out = Defense(EventHitSingleSS)
		} else {
			result, extra, out = Defense(EventHitSingle2B)
		}
	} else if roll == 6 || crit && roll == 5 {
		result, extra, out = Defense(EventHitSingleSS)
	} else if roll == 5 || crit && roll == 4 {
		result, extra, out = Defense(EventHitSingle3B)
	} else if roll == 4 || crit && roll == 3 {
		result, extra, out = Defense(EventHitSingle2B)
	} else if roll == 3 || crit && roll >= 1 {
		result, extra, out = Defense(EventHitSingle1B)
	} else {
		if p.HasTrait(TraitContact) {
			result, extra, out = EventHitDoubleAdv3, false, false
		} else {
			result, extra, out = EventHitSinglePlus, false, false
		}
	}

	return result, extra, out
}

package deadball

import (
	"fmt"
	"strings"

	"github.com/0xa1-red/phaseball/internal/dice"
	"github.com/google/uuid"
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
	ID      uuid.UUID  `db:"id"`
	Name    string     `db:"name"`
	Players [9]*Player `db:"-"`
	Index   int        `db:"-" json:"-" `
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

func GetPositionFromShort(short string) Position {
	switch short {
	case Catcher.Short:
		return Catcher
	case Pitcher.Short:
		return Pitcher
	case First.Short:
		return First
	case Second.Short:
		return Second
	case Third.Short:
		return Third
	case Shortstop.Short:
		return Shortstop
	case Left.Short:
		return Left
	case Center.Short:
		return Center
	case Right.Short:
		return Right
	}

	return Position{}
}

// Player represents a single player
type Player struct {
	Power        int       `db:"batter_pow"`
	Contact      int       `db:"batter_con"`
	Eye          int       `db:"batter_eye"`
	Speed        int       `db:"batter_spd"`
	Defense      int       `db:"batter_def"`
	Fastball     int       `db:"pitcher_fb"`
	Changeup     int       `db:"pitcher_ch"`
	Breaking     int       `db:"pitcher_bb"`
	Control      int       `db:"pitcher_ctl"`
	Batting      int       `db:"pitcher_bat"`
	Name         string    `db:"name"`
	Status       string    `db:"-" json:"-"`
	Position     Position  `db:"position"`
	BatterTarget uint8     `db:"-"`
	PitchDie     PitchDie  `db:"-"`
	Hand         string    `db:"hand"`
	ID           uuid.UUID `db:"id"`
	TeamID       uuid.UUID `db:"idteam"`
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

package deadball

import (
	"fmt"

	"hq.0xa1.red/axdx/phaseball/internal/dice"
)

type Event struct {
	Label string
	Long  string
	Extra string
}

var (
	EventHitSingle      Event = Event{Label: "SINGLE"}
	EventHitSinglePlus  Event = Event{Label: "SINGLE+", Long: "Single - Extra"}
	EventHitSingleError Event = Event{Label: "SINGLE_ERROR", Long: "Single - Error"}
	EventHitSingleAdv2  Event = Event{Label: "SINGLE_ADV_2", Long: "Single - Runners advance two"}
	EventHitSingle1B    Event = Event{Label: "SINGLE_DEF_1B"}
	EventHitSingle2B    Event = Event{Label: "SINGLE_DEF_2B"}
	EventHitSingle3B    Event = Event{Label: "SINGLE_DEF_3B"}
	EventHitSingleSS    Event = Event{Label: "SINGLE_DEF_SS"}

	EventHitDouble      Event = Event{Label: "DOUBLE"}
	EventHitDoubleError Event = Event{Label: "DOUBLE_ERROR"}
	EventHitDoubleAdv3  Event = Event{Label: "DOUBLE_ADV_3"}
	EventHitDoubleLF    Event = Event{Label: "DOUBLE_DEF_LF"}
	EventHitDoubleCF    Event = Event{Label: "DOUBLE_DEF_CF"}
	EventHitDoubleRF    Event = Event{Label: "DOUBLE_DEF_RF"}

	EventHitTriple      Event = Event{Label: "TRIPLE"}
	EventHitTripleError Event = Event{Label: "TRIPLE_ERROR"}
	EventHitTripleRF    Event = Event{Label: "TRIPLE_DEF_RF"}
	EventHitTripleCF    Event = Event{Label: "TRIPLE_DEF_CF"}

	EventHitHomeRun Event = Event{Label: "HOME_RUN", Long: "Home run", Extra: ""}

	EventHitProductiveOut  Event = Event{Label: "PRODUCTIVE_OUT", Long: "", Extra: ""}
	EventHitDoublePlay     Event = Event{Label: "DOUBLE_PLAY", Long: "", Extra: ""}
	EventHitFieldersChoice Event = Event{Label: "FIELDERS_CHOICE", Long: "", Extra: ""}

	EventWalk        Event = Event{Label: "WALK", Long: "Walk", Extra: ""}
	EventOut         Event = Event{Label: "OUT", Long: "Out", Extra: ""}
	EventCrit        Event = Event{Label: "CRITICAL_HIT", Long: "Critical hit", Extra: ""}
	EventError       Event = Event{Label: "ERROR", Long: "Error", Extra: ""}
	EventHit         Event = Event{Label: "HIT", Long: "Hit", Extra: ""}
	EventProdOut     Event = Event{Label: "PRODUCTIVE_OUT", Long: "", Extra: ""}
	EventPossibleDbl Event = Event{Label: "POSSIBLE_DOUBLE", Long: "", Extra: ""}

	EventOutK  Event = Event{Label: EventOut.Label, Long: "Strikeout", Extra: "K"}
	EventOutG3 Event = Event{Label: EventOut.Label, Long: "Groundout to first", Extra: "G-3"}
	EventOut43 Event = Event{Label: EventOut.Label, Long: "Groundout to second", Extra: "4-3"}
	EventOut53 Event = Event{Label: EventOut.Label, Long: "Groundout to third", Extra: "5-3"}
	EventOut63 Event = Event{Label: EventOut.Label, Long: "Groundout to short", Extra: "6-3"}
	EventOutF7 Event = Event{Label: EventOut.Label, Long: "Flyout to left field", Extra: "F-7"}
	EventOutF8 Event = Event{Label: EventOut.Label, Long: "Flyout to center field", Extra: "F-8"}
	EventOutF9 Event = Event{Label: EventOut.Label, Long: "Flyout to right field", Extra: "F-9"}
)

func (e Event) GetLong() string {
	if e.Long != "" {
		return e.Long
	}

	return string(e.Label)
}

func Hit(swing int, crit bool) (result Event, extra bool, out bool) {
	roll := dice.Roll(20, 1, 0)

	return hit(swing, crit, roll)
}

func hit(swing int, crit bool, roll int) (Event, bool, bool) {
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
		result, extra, out = EventHitSinglePlus, false, false
	}

	return result, extra, out
}

func Defense(hit Event) (result Event, extra bool, out bool) {
	roll := dice.Roll(12, 1, 0)

	return defense(hit, roll)
}

func defense(hit Event, roll int) (Event, bool, bool) {
	var result Event
	var out bool
	var extra bool
	if roll == 12 {
		result = hit
		out = true
	} else if roll >= 10 {
		switch hit {
		case EventHitDoubleCF, EventHitDoubleLF, EventHitDoubleRF:
			result = EventHitSingleAdv2
		case EventHitTripleCF, EventHitTripleRF:
			result = EventHitDoubleAdv3
		default:
			result = hit
		}
	} else if roll >= 3 {
		result = hit
	} else {
		result = hit
		extra = true
	}
	return result, out, extra
}

func Out(swing int) Event {
	events := map[string]Event{
		"0": EventOutK,
		"1": EventOutK,
		"2": EventOutK,
		"3": EventOutG3,
		"4": EventOut43,
		"5": EventOut53,
		"6": EventOut63,
		"7": EventOutF7,
		"8": EventOutF8,
		"9": EventOutF9,
	}
	s := fmt.Sprintf("%d", swing)
	digit := s[len(s)-1:]

	if res, ok := events[digit]; ok {
		return res
	}

	return Event{}
}

func IsOutInfield(e Event) bool {
	if e == EventOutG3 ||
		e == EventOut43 ||
		e == EventOut53 ||
		e == EventOut63 {
		return true
	}

	return false
}

func IsOutOutfield(e Event) bool {
	if e == EventOutF7 ||
		e == EventOutF8 ||
		e == EventOutF9 {
		return true
	}

	return false
}

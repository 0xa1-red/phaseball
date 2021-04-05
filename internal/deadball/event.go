package deadball

import (
	"fmt"

	"hq.0xa1.red/axdx/phaseball/internal/dice"
)

type ExtendedEvent struct {
	Event Event
	Long  string
	Extra string
}

type Event string

const (
	EventLogWalk Event = "WALK"

	EventHitSinglePlus  Event = "SINGLE+"
	EventHitSingleAdv2  Event = "SINGLE_ADV_2"
	EventHitSingle1B    Event = "SINGLE_DEF_1B"
	EventHitSingle2B    Event = "SINGLE_DEF_2B"
	EventHitSingle3B    Event = "SINGLE_DEF_3B"
	EventHitSingleSS    Event = "SINGLE_DEF_SS"
	EventHitSingleError Event = "SINGLE_ERROR"
	EventHitSingle      Event = "SINGLE"

	EventHitDoubleAdv3  Event = "DOUBLE_ADV_3"
	EventHitDoubleLF    Event = "DOUBLE_DEF_LF"
	EventHitDoubleCF    Event = "DOUBLE_DEF_CF"
	EventHitDoubleRF    Event = "DOUBLE_DEF_RF"
	EventHitDoubleError Event = "DOUBLE_ERROR"
	EventHitDouble      Event = "DOUBLE"

	EventHitTripleRF    Event = "TRIPLE_DEF_RF"
	EventHitTripleCF    Event = "TRIPLE_DEF_CF"
	EventHitTripleError Event = "TRIPLE_ERROR"
	EventHitTriple      Event = "TRIPLE"

	EventHitHomeRun Event = "HOME_RUN"
	EventHitOut     Event = "OUT"

	EventHitProductiveOut  Event = "PRODUCTIVE_OUT"
	EventHitDoublePlay     Event = "DOUBLE_PLAY"
	EventHitFieldersChoice Event = "FIELDERS_CHOICE"

	EventOutK  Event = "K"
	EventOutG3 Event = "G-3"
	EventOut43 Event = "4-3"
	EventOut53 Event = "5-3"
	EventOut63 Event = "6-3"
	EventOutF7 Event = "F-7"
	EventOutF8 Event = "F-8"
	EventOutF9 Event = "F-9"
)

var ExtendedEventMapping = map[Event]ExtendedEvent{
	EventHitHomeRun: {Event: EventHitHomeRun, Long: "Home run", Extra: ""},
	EventOutK:       {Event: EventHitOut, Long: "Strikeout", Extra: "K"},
	EventOutG3:      {Event: EventHitOut, Long: "Groundout to first", Extra: "G-3"},
	EventOut43:      {Event: EventHitOut, Long: "Groundout to second", Extra: "4-3"},
	EventOut53:      {Event: EventHitOut, Long: "Groundout to third", Extra: "5-3"},
	EventOut63:      {Event: EventHitOut, Long: "Groundout to short", Extra: "6-3"},
	EventOutF7:      {Event: EventHitOut, Long: "Flyout to left field", Extra: "F-7"},
	EventOutF8:      {Event: EventHitOut, Long: "Flyout to center field", Extra: "F-8"},
	EventOutF9:      {Event: EventHitOut, Long: "Flyout to right field", Extra: "F-9"},
}

func (e ExtendedEvent) GetLong() string {
	if e.Long != "" {
		return e.Long
	}

	return string(e.Event)
}

func (e Event) Short() string {
	return string(e)
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

	return ""
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

package deadball

// Event types
const (
	EventCrit int = iota
	EventHit
	EventWalk
	EventProdOut
	EventPossibleDbl
)

// Event type strings
const (
	EventOutStr         string = "Out"
	EventErrorStr       string = "Error"
	EventCritStr        string = "Critical hit"
	EventHitStr         string = "Hit"
	EventWalkStr        string = "Walk"
	EventProdOutStr     string = "Productive out"
	EventPossibleDblStr string = "Possible double"
)

// Team keys
const (
	TeamAway string = "away"
	TeamHome string = "home"
)

// Base names
const (
	BaseFirst  string = "First"
	BaseSecond string = "Second"
	BaseThird  string = "Third"
	BaseHome   string = "Home"
)

// Player status names
const (
	StatusBase    string = "On base"
	StatusOnDeck  string = "On deck"
	StatusWaiting string = "Waiting"
	StatusOut     string = "Out"
)

// Positions
const (
	PositionPitcher byte = iota
	PositionCatcher
	PositionFirst
	PositionSecond
	PositionThird
	PositionShortstop
	PositionLeft
	PositionCenter
	PositionRight
)

type PitchDie string

const (
	PitchNone   PitchDie = "NONE"
	PitchAddD12 PitchDie = "D12"
	PitchAddD8  PitchDie = "D8"
	PitchAddD4  PitchDie = "D4"
	PitchSubD4  PitchDie = "-D4"
)

const (
	HalfTop    string = "Top"
	HalfBottom string = "Bottom"
)

package deadball

import (
	"fmt"
	"strings"
	"testing"
)

func TestGame_Score(t *testing.T) {
	game := Game{
		Turns: []*Turn{
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
			{
				Top:    &Inning{Runs: 1},
				Bottom: &Inning{Runs: 0},
			},
		},
	}

	s := game.Score()

	if expected, actual := uint8(9), s[TeamAway]; expected != actual {
		t.Fatalf("Fail: expected away score %d, got %d", expected, actual)
	}

	if expected, actual := uint8(0), s[TeamHome]; expected != actual {
		t.Fatalf("Fail: expected home score %d, got %d", expected, actual)
	}
}

func TestInning_PossibleDouble(t *testing.T) {
	tests := []struct {
		outs          uint8
		label         string
		swing         int
		expectedEvent Event
		expectedOuts  uint8
	}{
		{
			outs:          1,
			label:         "double play",
			swing:         75,
			expectedEvent: Event{Label: EventHitDoublePlay.Label, Extra: "4-3"},
			expectedOuts:  2,
		},
		{
			outs:          1,
			label:         "fielders choice",
			swing:         66,
			expectedEvent: Event{Label: EventHitFieldersChoice.Label, Extra: "4-3"},
			expectedOuts:  1,
		},
		{
			outs:          1,
			label:         strings.ToLower(EventOut43.Long),
			swing:         51,
			expectedEvent: Event{Label: EventOut.Label, Extra: "4-3"},
			expectedOuts:  1,
		},
	}

	for i, tt := range tests {
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, tt.label)
			{
				GetDiamond().Reset()
				team.NewTurn(true)
				i := NewInning(&team, &Player{Name: "Peter Test"}, 1, HalfTop)
				i.Outs = tt.outs
				i.Diamond.Bases[0].Player = team.Players[0]

				actualEvent := i.PossibleDouble(tt.swing, EventOut43, team.Players[1])

				if expected, actual := tt.expectedEvent.Label, actualEvent.Label; expected != actual {
					t.Fatalf("Fail: expected event to be %s, got %s", expected, actual)
				}

				if expected, actual := tt.expectedEvent.Extra, actualEvent.Extra; expected != actual {
					t.Fatalf("Fail: expected extra to be %s, got %s", expected, actual)
				}

				if i.Outs != tt.expectedOuts {
					t.Fatalf("Fail: expected %d outs, got %d", tt.expectedOuts, i.Outs)
				}

				t.Log("Pass")
			}
		}

		t.Run(tt.label, tf)
	}
}

func TestInning_ProductiveOut(t *testing.T) {
	tests := []struct {
		label         string
		outEvent      Event
		p2            *Player
		p3            *Player
		expectedRuns  []*Player
		expectedEvent Event
	}{
		{
			label:         "no bases",
			outEvent:      EventOutF7,
			p2:            nil,
			p3:            nil,
			expectedRuns:  make([]*Player, 0),
			expectedEvent: EventOutF7,
		},
		{
			label:         "second base",
			outEvent:      EventOutF7,
			p2:            team.Players[0],
			p3:            nil,
			expectedRuns:  make([]*Player, 0),
			expectedEvent: Event{Label: EventHitProductiveOut.Label, Extra: "F-7"},
		},
		{
			label:         "second and third base",
			outEvent:      EventOutF7,
			p2:            team.Players[1],
			p3:            team.Players[0],
			expectedRuns:  []*Player{team.Players[0]},
			expectedEvent: Event{Label: EventHitProductiveOut.Label, Extra: "F-7"},
		},
		{
			label:         "third base",
			outEvent:      EventOutF7,
			p2:            nil,
			p3:            team.Players[0],
			expectedRuns:  []*Player{team.Players[0]},
			expectedEvent: Event{Label: EventHitProductiveOut.Label, Extra: "F-7"},
		},
	}

	for i, tt := range tests {
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, tt.label)
			{
				GetDiamond().Reset()
				team.NewTurn(true)
				i := NewInning(&team, &Player{Name: "Peter Test"}, 1, HalfTop)
				i.Outs = 1
				i.Diamond.Bases[1].Player = tt.p2
				i.Diamond.Bases[2].Player = tt.p3

				actualEvent, actualRunners := i.ProductiveOut(60, tt.outEvent)

				if expected, actual := tt.expectedEvent.Label, actualEvent.Label; expected != actual {
					t.Fatalf("Fail: expected event to be %s, got %s", expected, actual)
				}

				if expected, actual := tt.expectedEvent.Extra, actualEvent.Extra; expected != actual {
					t.Fatalf("Fail: expected extra to be %s, got %s", expected, actual)
				}

				if e, a := len(tt.expectedRuns), len(actualRunners); e != a {
					t.Fatalf("Fail: expected %d runners, got %d", e, a)
				}

				if len(tt.expectedRuns) > 0 && len(actualRunners) > 0 {
					if e, a := tt.expectedRuns[0], actualRunners[0]; e.Name != a.Name {
						t.Fatalf("Fail: expected %s to run home, got %s", e.Name, a.Name)
					}
				}

				t.Log("Pass")
			}
		}

		t.Run(tt.label, tf)
	}
}

func TestHomeRun(t *testing.T) {
	GetDiamond().Reset()
	team.NewTurn(true)

	inning := Inning{
		Outs:    0,
		Runs:    0,
		Hits:    0,
		Diamond: GetDiamond(),
	}

	tests := []struct {
		label         string
		description   string
		sequence      []int
		expectedBases []*Player
		expectedRuns  []*Player
	}{
		{
			label:       "bases_empty",
			description: "Bases are empty and the batter hits a HR",
			sequence:    []int{4},
			expectedBases: []*Player{
				nil,
				nil,
				nil,
				nil,
			},
			expectedRuns: []*Player{
				team.Players[0],
			},
		},
		{
			label:       "bases_scarce",
			description: "First and third are empty and the batter hits a HR",
			sequence:    []int{3, 1, 4},
			expectedBases: []*Player{
				nil,
				nil,
				nil,
				nil,
			},
			expectedRuns: []*Player{
				team.Players[0],
				team.Players[1],
				team.Players[2],
			},
		},
		{
			label:       "bases_loaded",
			description: "Bases are loaded and the 4th batter hits a HR",
			sequence:    []int{1, 1, 1, 4},
			expectedBases: []*Player{
				nil,
				nil,
				nil,
				nil,
			},
			expectedRuns: []*Player{
				team.Players[0],
				team.Players[1],
				team.Players[2],
				team.Players[3],
			},
		},
	}

	for i, tt := range tests {
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, tt.description)
			{
				GetDiamond().Reset()
				team.NewTurn(true)
				var actualRuns []*Player
				for _, step := range tt.sequence {
					actualRuns = append(actualRuns, inning.Diamond.Advance(team.AtBat(), step)...)
				}

				if len(tt.expectedRuns) != len(actualRuns) {
					t.Fatalf("Fail: expected %d runs, got %d", len(tt.expectedRuns), len(actualRuns))
				}

				for i, expectedRunner := range tt.expectedRuns {
					if actualRunner := actualRuns[i]; expectedRunner.Name != actualRunner.Name {
						t.Fatalf("Fail: expected %s to be runner #%d but they aren't", expectedRunner.Name, i+1)
					}
				}

				for i := range inning.Diamond.Bases {
					if expected, actual := tt.expectedBases[i], inning.Diamond.Bases[i].Player; expected != actual {
						t.Fatalf("Fail: expected %+v on base #%d, got %+v", expected, i+1, actual)
					}
				}
				t.Log("Pass")
			}
		}

		t.Run(tt.label, tf)
	}
}

func TestSinglesGetRun(t *testing.T) {
	GetDiamond().Reset()
	team.NewTurn(true)

	inning := Inning{
		Outs:    0,
		Runs:    0,
		Hits:    0,
		Diamond: GetDiamond(),
	}

	expectedBases := []*Player{
		team.Players[3],
		team.Players[2],
		team.Players[1],
		nil,
	}

	var expectedRuns []*Player = []*Player{
		team.Players[0],
	}
	var actualRuns []*Player
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[0], 1)...)
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[1], 1)...)
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[2], 1)...)
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[3], 1)...)

	t.Log("Given there were 4 singles, we expect one run and the last 3 players on bases")
	{
		if len(expectedRuns) != len(actualRuns) {
			t.Fatalf("Fail: expected %d runs, got %d", len(expectedRuns), len(actualRuns))
		}

		for i, expectedRunner := range expectedRuns {
			if actualRunner := actualRuns[i]; expectedRunner.Name != actualRunner.Name {
				t.Fatalf("Fail: expected %s to be runner #%d but they aren't", expectedRunner.Name, i+1)
			}
		}

		for i := range inning.Diamond.Bases {
			if expected, actual := expectedBases[i], inning.Diamond.Bases[i].Player; expected != actual {
				t.Fatalf("Fail: expected %+v on base #%d, got %+v", expected, i+1, actual)
			}
		}
		t.Log("Pass")
	}
}

func TestTwoSinglesAndDouble(t *testing.T) {
	GetDiamond().Reset()
	team.NewTurn(true)

	inning := Inning{
		Outs:    0,
		Runs:    0,
		Hits:    0,
		Diamond: GetDiamond(),
	}

	expectedBases := []*Player{
		nil,
		team.Players[2],
		team.Players[1],
		nil,
	}

	var expectedRuns []*Player = []*Player{
		team.Players[0],
	}
	var actualRuns []*Player
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[0], 1)...)
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[1], 1)...)
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[2], 2)...)

	t.Log("Given there were 2 singles and 1 double, we expect one run and the last 2 players on bases")
	{
		if len(expectedRuns) != len(actualRuns) {
			t.Fatalf("Fail: expected %d runs, got %d", len(expectedRuns), len(actualRuns))
		}

		for i, expectedRunner := range expectedRuns {
			if actualRunner := actualRuns[i]; expectedRunner.Name != actualRunner.Name {
				t.Fatalf("Fail: expected %s to be runner #%d but they aren't", expectedRunner.Name, i+1)
			}
		}

		for i := range inning.Diamond.Bases {
			if expected, actual := expectedBases[i], inning.Diamond.Bases[i].Player; expected != actual {
				t.Fatalf("Fail: expected %+v on base #%d, got %+v", expected, i+1, actual)
			}
		}
		t.Log("Pass")
	}
}

func TestOneSingleAndTriple(t *testing.T) {
	GetDiamond().Reset()
	team.NewTurn(true)

	inning := Inning{
		Outs:    0,
		Runs:    0,
		Hits:    0,
		Diamond: GetDiamond(),
	}

	expectedBases := []*Player{
		nil,
		nil,
		team.Players[1],
		nil,
	}

	var expectedRuns []*Player = []*Player{
		team.Players[0],
	}
	var actualRuns []*Player
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[0], 1)...)
	actualRuns = append(actualRuns, inning.Diamond.Advance(team.Players[1], 3)...)

	t.Log("Given there were 1 single and 1 triple, we expect one run and the last player on bases")
	{
		if len(expectedRuns) != len(actualRuns) {
			t.Fatalf("Fail: expected %d runs, got %d", len(expectedRuns), len(actualRuns))
		}

		for i, expectedRunner := range expectedRuns {
			if actualRunner := actualRuns[i]; expectedRunner.Name != actualRunner.Name {
				t.Fatalf("Fail: expected %s to be runner #%d but they aren't", expectedRunner.Name, i+1)
			}
		}

		for i := range inning.Diamond.Bases {
			if expected, actual := expectedBases[i], inning.Diamond.Bases[i].Player; expected != actual {
				t.Fatalf("Fail: expected %+v on base #%d, got %+v", expected, i+1, actual)
			}
		}
		t.Log("Pass")
	}
}

func TestSwingEvent(t *testing.T) {
	tests := []struct {
		swing         int
		bt            int
		expectedEvent Event
	}{
		{
			swing:         71,
			expectedEvent: EventPossibleDbl,
		},
		{
			swing:         60,
			bt:            50,
			expectedEvent: EventProdOut,
		},
		{
			swing:         52,
			bt:            50,
			expectedEvent: EventWalk,
		},
		{
			swing:         10,
			bt:            10,
			expectedEvent: EventHit,
		},
		{
			swing:         1,
			bt:            10,
			expectedEvent: EventCrit,
		},
	}

	for i, tt := range tests {
		label := fmt.Sprintf("Case %d - Swing: %d - BT: %d", i+1, tt.swing, tt.bt)
		tf := func(t *testing.T) {
			t.Log(label)

			if expected, actual := tt.expectedEvent, swingEvent(tt.swing, tt.bt); expected != actual {
				t.Fatalf("Fail: expected %s, got %s", expected, actual)
			}

			t.Log("Pass")
		}

		t.Run(label, tf)
	}
}

func TestDiamond_String(t *testing.T) {
	d := newDiamond()
	d.Bases[1].Player = team.Players[0]

	expected := "\tFirst : empty\n\tSecond : Anna Test\n\tThird : empty\n\tHome : empty\n"

	if actual := d.String(); expected != actual {
		t.Fatalf("Fail: expected:\n%s\ngot:\n%s", expected, actual)
	}

	t.Log("Pass")
}

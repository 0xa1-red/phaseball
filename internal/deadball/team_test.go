package deadball

import (
	"fmt"
	"testing"
)

func TestPlayer_Pitch(t *testing.T) {
	tests := []struct {
		label       string
		batter      *Player
		pitcher     *Player
		expectedDie PitchDie
	}{
		{
			label:       "lhp v lhp -4",
			batter:      &Player{Hand: HandLeftie},
			pitcher:     &Player{PitchDie: PitchSubD4, Hand: HandLeftie},
			expectedDie: PitchAddD4,
		},
		{
			label:       "lhp v lhp 4",
			batter:      &Player{Hand: HandLeftie},
			pitcher:     &Player{PitchDie: PitchAddD4, Hand: HandLeftie},
			expectedDie: PitchAddD8,
		},
		{
			label:       "lhp v lhp 8",
			batter:      &Player{Hand: HandLeftie},
			pitcher:     &Player{PitchDie: PitchAddD8, Hand: HandLeftie},
			expectedDie: PitchAddD12,
		},
		{
			label:       "lhp v lhp 12",
			batter:      &Player{Hand: HandLeftie},
			pitcher:     &Player{PitchDie: PitchAddD12, Hand: HandLeftie},
			expectedDie: PitchAddD12,
		},
		{
			label:       "rhp v lhp -4",
			batter:      &Player{Hand: HandRightie},
			pitcher:     &Player{PitchDie: PitchSubD4, Hand: HandLeftie},
			expectedDie: PitchSubD4,
		},
		{
			label:       "rhp v lhp 4",
			batter:      &Player{Hand: HandRightie},
			pitcher:     &Player{PitchDie: PitchAddD4, Hand: HandLeftie},
			expectedDie: PitchAddD4,
		},
		{
			label:       "rhp v lhp 8",
			batter:      &Player{Hand: HandRightie},
			pitcher:     &Player{PitchDie: PitchAddD8, Hand: HandLeftie},
			expectedDie: PitchAddD8,
		},
		{
			label:       "rhp v lhp 12",
			batter:      &Player{Hand: HandRightie},
			pitcher:     &Player{PitchDie: PitchAddD12, Hand: HandLeftie},
			expectedDie: PitchAddD12,
		},
	}
	for i, tt := range tests {
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, tt.label)

			actualDie, _ := tt.pitcher.Pitch(tt.batter.Hand)

			if actualDie != tt.expectedDie {
				t.Fatalf("Fail: expected pitch die %s, got %s", tt.expectedDie, actualDie)
			}

			t.Log("Pass")
		}

		t.Run(tt.label, tf)
	}
}

func TestTeam_OnDeck(t *testing.T) {
	tests := []struct {
		index          int
		waiting        bool
		expectedPlayer *Player
	}{
		{
			index:          1,
			waiting:        true,
			expectedPlayer: team.Players[1],
		},
		{
			index:          8,
			waiting:        true,
			expectedPlayer: team.Players[8],
		},
		{
			index:          8,
			waiting:        false,
			expectedPlayer: team.Players[0],
		},
	}

	for i, tt := range tests {
		GetDiamond().Reset()
		team.NewTurn(true)
		tf := func(t *testing.T) {
			t.Logf("Case %d - %d\n", i+1, tt.index)

			team.Index = tt.index
			if !tt.waiting {
				team.Players[tt.index].Status = StatusOut
			}

			if expected, actual := tt.expectedPlayer.Name, team.OnDeck().Name; expected != actual {
				t.Fatalf("Fail: expected player at bat to be %s, got %s", expected, actual)
			}

			t.Log("Pass")
		}

		t.Run(fmt.Sprintf("index_%d", tt.index), tf)
	}
}

func TestTeam_Pitcher(t *testing.T) {
	if expected, actual := team.Players[0].Name, team.Pitcher().Name; expected != actual {
		t.Fatalf("Fail: expected pitcher to be %s, got %s", expected, actual)
	}

	team.Players[0].Position = Catcher

	var expected *Player = nil

	if actual := team.Pitcher(); expected != actual {
		t.Fatalf("Fail: expected pitcher to be nil, got %v", actual)
	}

	t.Log("Pass")
}

func TestTeam_Next(t *testing.T) {
	tests := []struct {
		index    int
		expected int
	}{
		{
			index:    1,
			expected: 2,
		},
		{
			index:    8,
			expected: 0,
		},
	}

	for i, tt := range tests {
		GetDiamond().Reset()
		team.NewTurn(true)
		tf := func(t *testing.T) {
			t.Logf("Case %d - %d\n", i+1, tt.index)

			team.Index = tt.index

			team.Next()

			if expected, actual := tt.expected, team.Index; expected != actual {
				t.Fatalf("Fail: expected index %d, got %d", expected, actual)
			}

			t.Log("Pass")
		}

		t.Run(fmt.Sprintf("index_%d", tt.index), tf)
	}
}

func TestTeam_String(t *testing.T) {
	team.NewTurn(true)
	expected := "\tAnna Test (On deck)\n\tBob Test (Waiting)\n\tClyde Test (Waiting)\n\tDoris Test (Waiting)\n\tElmer Test (Waiting)\n\tFrank Test (Waiting)\n\tGillian Test (Waiting)\n\tHelen Test (Waiting)\n\tIan Test (Waiting)"
	actual := team.String()

	if expected != actual {
		t.Fatalf("Fail: expected:\n%s\ngot:\n%s", expected, actual)
	}

	t.Log("Pass")
}

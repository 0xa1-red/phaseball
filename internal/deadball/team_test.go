package deadball

import "testing"

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

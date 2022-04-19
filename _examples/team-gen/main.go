package main

import (
	"fmt"

	"github.com/0xa1-red/phaseball/internal/dice"
)

func main() {
	positions := []string{
		"deadball.Pitcher",
		"deadball.Catcher",
		"deadball.First",
		"deadball.Second",
		"deadball.Third",
		"deadball.Shortstop",
		"deadball.Left",
		"deadball.Center",
		"deadball.Right",
	}

	names := []string{
		"Anna Home",
		"Bob Home",
		"Clyde Home",
		"Doris Home",
		"Elmer Home",
		"Frank Home",
		"Gillian Home",
		"Helen Home",
		"Ian Home",
	}

	for i := 0; i < 9; i++ {
		var batterTarget uint8
		var pitchDie string
		if positions[i] == "deadball.Pitcher" {
			batterTarget = uint8(dice.Roll(10, 1, 5))
			pitchDieRoll := dice.Roll(8, 1, 0)
			if pitchDieRoll > 6 {
				pitchDie = "deadball.PitchSubD4"
			} else if pitchDieRoll > 3 {
				pitchDie = "deadball.PitchAddD4"
			} else if pitchDieRoll > 1 {
				pitchDie = "deadball.PitchAddD8"
			} else {
				pitchDie = "deadball.PitchAddD12"
			}
		} else {
			batterTarget = uint8(dice.Roll(10, 2, 15))
			pitchDie = "deadball.PitchNone"
		}

		fmt.Printf(`{
	Name: "%s",
	Status: deadball.StatusWaiting,
	Position: %s,
	BatterTarget: %d,
	PitchDie: %s,
},
`, names[i], positions[i], batterTarget, pitchDie)
	}

}

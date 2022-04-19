package main

import (
	"fmt"
	"log"

	"github.com/0xa1-red/phaseball/internal/deadball"
)

func main() {
	away := deadball.Team{
		Name: "Away Avengers",
		Players: [9]*deadball.Player{
			{
				Name:         "Anna Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Pitcher,
				BatterTarget: 8,
				PitchDie:     deadball.PitchAddD12,
			},
			{
				Name:         "Bob Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Catcher,
				BatterTarget: 23,
				PitchDie:     deadball.PitchNone,
			},
			{
				Name:         "Clyde Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.First,
				BatterTarget: 28,
				PitchDie:     deadball.PitchNone,
			},
			{
				Name:         "Doris Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Second,
				BatterTarget: 25,
				PitchDie:     deadball.PitchNone,
			},
			{
				Name:         "Elmer Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Third,
				BatterTarget: 24,
				PitchDie:     deadball.PitchNone,
			},
			{
				Name:         "Frank Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Shortstop,
				BatterTarget: 30,
				PitchDie:     deadball.PitchNone,
			},
			{
				Name:         "Gillian Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Left,
				BatterTarget: 25,
				PitchDie:     deadball.PitchNone,
			},
			{
				Name:         "Helen Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Center,
				BatterTarget: 26,
				PitchDie:     deadball.PitchNone,
			},
			{
				Name:         "Ian Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Right,
				BatterTarget: 24,
				PitchDie:     deadball.PitchNone,
			},
		},
	}
	d := deadball.GetDiamond()

	inning := deadball.Inning{
		Team: &away,
		Pitcher: &deadball.Player{
			Name:         "Anna Away",
			Status:       deadball.StatusWaiting,
			Position:     deadball.Pitcher,
			BatterTarget: 8,
			PitchDie:     deadball.PitchAddD12,
		},
		Diamond: d,
	}

	fill := 3
	var runs uint8
	runs += inning.Diamond.Advance(inning.Team.AtBat(), 1)
	fmt.Println("")
	for i := 0; i < fill; i++ {
		if inning.Diamond.Bases[0].Load(nil) {
			runs++
		}
		fmt.Println("")
	}

	for _, base := range inning.Diamond.Bases {
		p := "empty"
		if base.Player != nil {
			p = base.Player.Name
		}
		fmt.Printf("%s : %s\n", base.Name, p)
	}

	log.Printf("Runs: %d\n", runs)
}

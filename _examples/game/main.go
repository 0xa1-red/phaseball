package main

import (
	"fmt"
	"os"

	"hq.0xa1.red/axdx/phaseball/internal/deadball"
)

func main() {
	away := deadball.Team{
		Name:  "Away Avengers",
		Index: 0,
		Players: [9]*deadball.Player{
			{
				Name:         "Anna Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Pitcher,
				BatterTarget: 8,
				PitchDie:     deadball.PitchAddD12,
				Hand:         deadball.HandRightie,
			},
			{
				Name:         "Bob Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Catcher,
				BatterTarget: 23,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandLeftie,
			},
			{
				Name:         "Clyde Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.First,
				BatterTarget: 28,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandRightie,
			},
			{
				Name:         "Doris Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Second,
				BatterTarget: 25,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandLeftie,
			},
			{
				Name:         "Elmer Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Third,
				BatterTarget: 24,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandRightie,
			},
			{
				Name:         "Frank Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Shortstop,
				BatterTarget: 30,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandRightie,
			},
			{
				Name:         "Gillian Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Left,
				BatterTarget: 25,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandLeftie,
			},
			{
				Name:         "Helen Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Center,
				BatterTarget: 26,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandLeftie,
			},
			{
				Name:         "Ian Home",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Right,
				BatterTarget: 24,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandRightie,
			},
		},
	}
	home := deadball.Team{
		Name:  "Home Heroes",
		Index: 0,
		Players: [9]*deadball.Player{
			{
				Name:         "Anna Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Pitcher,
				BatterTarget: 11,
				PitchDie:     deadball.PitchAddD8,
				Hand:         deadball.HandSwitch,
			},
			{
				Name:         "Bob Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Catcher,
				BatterTarget: 23,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandLeftie,
			},
			{
				Name:         "Clyde Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.First,
				BatterTarget: 28,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandLeftie,
			},
			{
				Name:         "Doris Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Second,
				BatterTarget: 25,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandSwitch,
			},
			{
				Name:         "Elmer Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Third,
				BatterTarget: 22,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandLeftie,
			},
			{
				Name:         "Frank Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Shortstop,
				BatterTarget: 20,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandRightie,
			},
			{
				Name:         "Gillian Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Left,
				BatterTarget: 34,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandLeftie,
			},
			{
				Name:         "Helen Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Center,
				BatterTarget: 26,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandRightie,
			},
			{
				Name:         "Ian Away",
				Status:       deadball.StatusWaiting,
				Position:     deadball.Right,
				BatterTarget: 23,
				PitchDie:     deadball.PitchNone,
				Hand:         deadball.HandRightie,
			},
		},
	}

	game := deadball.NewGame(away, home)

	game.Run()

	// inning := deadball.Inning{
	// 	Outs:    0,
	// 	Runs:    0,
	// 	Hits:    0,
	// 	Team:    team,
	// 	Diamond: deadball.GetDiamond(),
	// }

	// inning.Run()

	// for _, base := range inning.Diamond.Bases {
	// 	name := "empty"
	// 	if base.Player != nil {
	// 		name = base.Player.Name
	// 	}
	// 	fmt.Printf("%s : %s\n", base.Name, name)
	// }
	// fmt.Printf("Runs: %d\n", inning.Runs)
	if err := os.WriteFile("game.json", []byte(game.Log.String()), 0655); err != nil {
		fmt.Printf("error saving game log: %v\n", err)
	}
}

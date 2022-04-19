package main

import (
	"fmt"
	"os"

	"github.com/0xa1-red/phaseball/internal/deadball"
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

	game := deadball.New(away, home)

	game.Run()

	if err := os.WriteFile("game.json", []byte(game.Log.String()), 0655); err != nil {
		fmt.Printf("error saving game log: %v\n", err)
	}
}

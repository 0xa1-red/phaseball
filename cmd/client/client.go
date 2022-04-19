package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/0xa1-red/phaseball/internal/deadball"
	"github.com/0xa1-red/phaseball/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
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

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service.NewMatchSimulatorClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tt := fmt.Sprintf("%d", time.Now().UnixNano())

	jAway, err := json.Marshal(away.Players)
	if err != nil {
		panic(err)
	}
	jHome, err := json.Marshal(home.Players)
	if err != nil {
		panic(err)
	}

	homeSheet := service.TeamSheet{
		Name:    home.Name,
		Players: string(jHome),
	}

	awaySheet := service.TeamSheet{
		Name:    away.Name,
		Players: string(jAway),
	}

	r, err := c.SimulateGame(ctx, &service.GameRequest{
		GameID:    uuid.New().String(),
		Timestamp: tt,
		Home:      &homeSheet,
		Away:      &awaySheet,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	if err := os.WriteFile(fmt.Sprintf("%s.json", r.GetGameID()), []byte(r.GetLog()), 0755); err != nil {
		panic(err)
	}

}

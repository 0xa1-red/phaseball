package main

import (
	"flag"
	"log"
	"os"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/database"
	"github.com/0xa1-red/phaseball/internal/deadball"
	"github.com/0xa1-red/phaseball/internal/logger"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	configPath string
	homeID     string
	awayID     string
	pbp        bool
)

func main() {
	flag.BoolVar(&pbp, "play-by-play", false, "Whether to print the game in play-by-play format or not (press enter for next play)")
	flag.StringVar(&configPath, "cfg", "./config.yml", "Configuration file in YAML format")
	flag.StringVar(&awayID, "away", "", "Away team ID")
	flag.StringVar(&homeID, "home", "", "Home team ID")
	flag.Parse()

	if awayID == "" || homeID == "" {
		flag.Usage()
		os.Exit(1)
	}

	awayUUID := uuid.MustParse(awayID)
	homeUUID := uuid.MustParse(homeID)

	if err := config.Init(configPath); err != nil {
		panic(err)
	}

	db, err := database.Connection()
	if err != nil {
		log.Println("conn")
		panic(err)
	}

	away, err := db.GetTeam(awayUUID)
	if err != nil {
		log.Println("home")
		panic(err)
	}

	home, err := db.GetTeam(homeUUID)
	if err != nil {
		log.Println("home")
		panic(err)
	}

	opts := make([]deadball.GameOption, 0)
	isStdoutLog := config.Get().GameLog.Kind == logger.KindStdout || config.Get().GameLog.Kind == logger.KindPretty
	if pbp && isStdoutLog {
		opts = append(opts, deadball.WithPlayByPlay())
	}
	game := deadball.New(away, home, opts...)
	gameModel := database.Game{
		ID: game.ID,
		Teams: database.TeamList{
			Away: game.Teams[deadball.TeamAway].ID,
			Home: game.Teams[deadball.TeamHome].ID,
		},
	}
	if err := db.SaveGame(gameModel); err != nil {
		panic(err)
	}

	game.Run()

	game.NewLog.Close()
}

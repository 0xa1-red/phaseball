package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/database"
	"github.com/0xa1-red/phaseball/internal/deadball"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	configPath string
	homeID     string
	awayID     string
)

func main() {
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

	game := deadball.New(away, home)
	if err := db.SaveGame(game); err != nil {
		panic(err)
	}

	game.Run()

	if err := os.WriteFile("game.json", []byte(game.Log.String()), 0655); err != nil {
		fmt.Printf("error saving game log: %v\n", err)
	}
}

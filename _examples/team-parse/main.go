package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strconv"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/database"
	"github.com/0xa1-red/phaseball/internal/deadball"
)

var (
	csvPath    string
	configPath string
	teamName   string
)

func main() {
	flag.StringVar(&teamName, "name", "", "The team name")
	flag.StringVar(&csvPath, "team", "./team.csv", "Team file in CSV format")
	flag.StringVar(&configPath, "cfg", "./config.yml", "Configuration file in YAML format")
	flag.Parse()

	if teamName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := config.Init(configPath); err != nil {
		panic(err)
	}

	db, err := database.Connection()
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	fp, err := os.OpenFile(csvPath, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	decoder := csv.NewReader(fp)
	decoder.Comma = ';'
	records, err := decoder.ReadAll()
	if err != nil {
		panic(err)
	}

	players := [9]*deadball.Player{}
	for i, record := range records {

		var (
			pow, con, eye, spd, def int
			fb, ch, bb, cn, bt      int
		)
		if record[1] == "P" {
			fb, _ = strconv.Atoi(record[4])
			ch, _ = strconv.Atoi(record[5])
			bb, _ = strconv.Atoi(record[6])
			cn, _ = strconv.Atoi(record[7])
			bt, _ = strconv.Atoi(record[8])
		} else {
			pow, _ = strconv.Atoi(record[4])
			con, _ = strconv.Atoi(record[5])
			eye, _ = strconv.Atoi(record[6])
			spd, _ = strconv.Atoi(record[7])
			def, _ = strconv.Atoi(record[8])
		}

		players[i] = &deadball.Player{
			Name:     record[0],
			Position: getPos(record[1]),
			Hand:     getHand(record[3]),
			Power:    pow,
			Contact:  con,
			Eye:      eye,
			Speed:    spd,
			Defense:  def,
			Fastball: fb,
			Changeup: ch,
			Breaking: bb,
			Control:  cn,
			Batting:  bt,
		}
	}

	team := deadball.Team{
		Name:    teamName,
		Players: players,
	}

	if err := db.SaveTeam(team); err != nil {
		panic(err)
	}
}

func getHand(str string) string {
	switch str {
	case "L":
		return deadball.HandLeftie
	case "R":
		return deadball.HandRightie
	case "S":
		return deadball.HandSwitch
	}
	return ""
}

func getPos(str string) deadball.Position {
	switch str {
	case "P":
		return deadball.Pitcher
	case "C":
		return deadball.Catcher
	case "1B":
		return deadball.First
	case "2B":
		return deadball.Second
	case "3B":
		return deadball.Third
	case "SS":
		return deadball.Shortstop
	case "LF":
		return deadball.Left
	case "CF":
		return deadball.Center
	case "RF":
		return deadball.Right
	}
	return deadball.Position{}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/0xa1-red/phaseball/internal/deadball"
)

var file string

func main() {
	flag.StringVar(&file, "file", "./game.json", "game log path")
	flag.Parse()

	f, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var game deadball.GameLog
	err = json.Unmarshal(f, &game)
	if err != nil {
		panic(err)
	}

	var half string = ""

	for _, entry := range game.Entries {
		if entry.Inning.Half != half {
			half = entry.Inning.Half
			fmt.Printf("Inning %d - %s; Pitcher: %s\n", entry.Inning.Number, half, entry.Pitcher.Name)
		}

		event := entry.Event.GetLong()
		if entry.Event.Extra != "" {
			event = fmt.Sprintf("%s %s", event, entry.Event.Extra)
		}

		fmt.Printf("\tBatter: %-20s | %s\n", entry.Batter.Name, event)
	}
}

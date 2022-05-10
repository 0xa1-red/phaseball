package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

var csvPath string

func main() {
	flag.StringVar(&csvPath, "team", "./team.csv", "Team file in CSV format")
	flag.Parse()
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

	for _, record := range records {
		tpl := `{
	Name:         "%s",
	Status:       deadball.StatusWaiting,
	Position:     %s,
	Hand:         %s,
	Power:        %s,
	Contact:      %s,
	Eye:          %s,
	Speed:        %s,
	Defense:      %s,
}`
		if record[1] == "P" {
			tpl = `{
	Name:         "%s",
	Status:       deadball.StatusWaiting,
	Position:     %s,
	Hand:         %s,
	Fastball:     %s,
	Changeup:     %s,
	Breaking:     %s,
	Control:      %s,
	Batting:      %s,
}`
		}

		fmt.Printf(tpl+",\n",
			record[0],
			posString(record[1]),
			handString(record[3]),
			record[4],
			record[5],
			record[6],
			record[7],
			record[8],
		)
	}
}

func handString(str string) string {
	switch str {
	case "L":
		return "deadball.HandLeftie"
	case "R":
		return "deadball.HandRightie"
	case "S":
		return "deadball.HandSwitch"
	}
	return ""
}

func posString(str string) string {
	switch str {
	case "P":
		return "deadball.Pitcher"
	case "C":
		return "deadball.Catcher"
	case "1B":
		return "deadball.First"
	case "2B":
		return "deadball.Second"
	case "3B":
		return "deadball.Third"
	case "SS":
		return "deadball.Shortstop"
	case "LF":
		return "deadball.Left"
	case "CF":
		return "deadball.Center"
	case "RF":
		return "deadball.Right"
	}
	return ""
}

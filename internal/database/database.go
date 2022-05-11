package database

import (
	"fmt"
	"log"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/deadball"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var db *Conn

type Conn struct {
	*sqlx.DB
}

func Connection() (*Conn, error) {
	if db == nil {
		url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?ssl_mode=disable&options=--cluster=purple-moose-1962",
			config.Get().Database.User,
			config.Get().Database.Password,
			config.Get().Database.Host,
			config.Get().Database.Port,
			config.Get().Database.Db,
		)
		c, err := sqlx.Connect("postgres", url)
		if err != nil {
			log.Fatal(err)
		}

		db = &Conn{c}
	}

	return db, nil
}

func (c *Conn) SaveTeam(team deadball.Team) error {
	tx, err := c.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO teams (name) VALUES ($1)", team.Name); err != nil {
		tx.Rollback() // nolint
		return err
	}

	res := tx.QueryRow("SELECT id FROM teams WHERE name = $1", team.Name)
	var teamID string
	if err := res.Scan(&teamID); err != nil {
		tx.Rollback() // nolint
		return err
	}

	for _, player := range team.Players {
		_, err := tx.Exec(`INSERT INTO players (
			idteam,
			name,
			position,
			hand,
			batter_pow,
			batter_con,
			batter_eye,
			batter_spd,
			batter_def,
			pitcher_fb,
			pitcher_ch,
			pitcher_bb,
			pitcher_ctl,
			pitcher_bat
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
`, teamID, player.Name, player.Position.Short, player.Hand[0:1],
			player.Power, player.Contact, player.Eye, player.Speed, player.Defense,
			player.Fastball, player.Changeup, player.Breaking, player.Control, player.Batting)
		if err != nil {
			tx.Rollback() // nolint
			return err
		}
	}

	tx.Commit() // nolint

	return nil
}

func (c *Conn) GetTeam(id uuid.UUID) (deadball.Team, error) {
	res := c.QueryRowx("SELECT name FROM teams WHERE id = $1", id.String())
	if res.Err() != nil {
		return deadball.Team{}, res.Err()
	}

	var name string
	if err := res.Scan(&name); err != nil {
		return deadball.Team{}, err
	}

	t := deadball.Team{
		ID:   id,
		Name: name,
	}

	playerRows, err := c.Queryx("SELECT * FROM players WHERE idteam = $1", id.String())
	if err != nil {
		return deadball.Team{}, err
	}
	defer playerRows.Close()

	i := 0
	for playerRows.Next() {
		if playerRows.Err() != nil {
			return deadball.Team{}, playerRows.Err()
		}

		if err != nil {
			return deadball.Team{}, err
		}

		p := make(map[string]interface{})
		if err := playerRows.MapScan(p); err != nil {
			return deadball.Team{}, err
		}

		id, err := uuid.ParseBytes(p["id"].([]byte))
		if err != nil {
			return deadball.Team{}, err
		}

		idteam, err := uuid.ParseBytes(p["idteam"].([]byte))
		if err != nil {
			return deadball.Team{}, err
		}

		h, ok := p["hand"].([]byte)
		if !ok {
			return deadball.Team{}, fmt.Errorf("Invalid string assertion for hand")
		}
		hand := deadball.HandRightie
		switch string(h) {
		case "L":
			hand = deadball.HandLeftie
		case "S":
			hand = deadball.HandSwitch
		}

		pos, ok := p["position"].([]byte)
		if !ok {
			return deadball.Team{}, fmt.Errorf("Invalid string assertion for position")
		}

		player := &deadball.Player{
			ID:       id,
			TeamID:   idteam,
			Power:    int(p["batter_pow"].(int64)),
			Contact:  int(p["batter_con"].(int64)),
			Eye:      int(p["batter_eye"].(int64)),
			Speed:    int(p["batter_spd"].(int64)),
			Defense:  int(p["batter_def"].(int64)),
			Fastball: int(p["pitcher_fb"].(int64)),
			Changeup: int(p["pitcher_ch"].(int64)),
			Breaking: int(p["pitcher_bb"].(int64)),
			Control:  int(p["pitcher_ctl"].(int64)),
			Batting:  int(p["pitcher_bat"].(int64)),
			Name:     p["name"].(string),
			Position: deadball.GetPositionFromShort(string(pos)),
			Hand:     hand,
		}

		t.Players[i] = player
		i++
	}

	return t, nil
}

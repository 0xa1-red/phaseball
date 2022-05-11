package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/deadball"
	_ "github.com/lib/pq"
)

var db *Conn

type Conn struct {
	*sql.DB
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
		c, err := sql.Open("postgres", url)
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
	if res.Err() != nil {
		tx.Rollback() // nolint
		return err
	}

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

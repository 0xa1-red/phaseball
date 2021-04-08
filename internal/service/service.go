package service

import (
	context "context"
	"encoding/json"
	"log"
	"net"

	grpc "google.golang.org/grpc"
	"hq.0xa1.red/axdx/phaseball/internal/deadball"
)

var port = ":5051"

type Simulator struct {
	UnimplementedMatchSimulatorServer
}

func (s *Simulator) SimulateGame(ctx context.Context, g *GameRequest) (*MatchResults, error) {

	var awayPlayers [9]*deadball.Player
	err := json.Unmarshal([]byte(g.Away.Players), &awayPlayers)
	if err != nil {
		return &MatchResults{}, err
	}
	away := deadball.Team{
		Name:    g.Away.Name,
		Players: awayPlayers,
	}

	var homePlayers [9]*deadball.Player
	err = json.Unmarshal([]byte(g.Home.Players), &homePlayers)
	if err != nil {
		return &MatchResults{}, err
	}
	home := deadball.Team{
		Name:    g.Home.Name,
		Players: homePlayers,
	}

	game := deadball.NewGame(away, home)
	game.Run()

	return &MatchResults{
		Status: 1,
		GameID: g.GetGameID(),
		Log:    game.Log.String(),
	}, nil
}

func Start(done chan struct{}) *grpc.Server {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterMatchSimulatorServer(s, &Simulator{})
	go func(d chan struct{}) {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		d <- struct{}{}
	}(done)

	return s
}

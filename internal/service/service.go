package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/websocket"
)

type Server struct {
	*http.Server

	errors chan error
}

func New() *Server {
	m := mux.NewRouter()

	m.HandleFunc("/replay/{game_id}", handleGameReplay)
	m.HandleFunc("/game/{game_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID, err := uuid.Parse(vars["game_id"])
		if err != nil {
			http.Error(w, "Missing Game ID: GET /replay/<game_id>", http.StatusNotFound)
			return
		}

		db, err := database.Connection()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error establishing database connection", http.StatusInternalServerError)
			return
		}

		game, err := db.GetGameLog(gameID)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error retrieving game log", http.StatusInternalServerError)
			return
		}

		log.Println(game)
		w.Write([]byte("Shrug")) // nolint
	})

	s := &http.Server{
		Addr:    config.Get().Service.Address,
		Handler: m,
	}

	ss := &Server{Server: s}

	go ss.Start()

	return &Server{Server: s}
}

func (s *Server) Errors() chan error {
	return s.errors
}

func (s *Server) Start() {
	fmt.Printf("Listening on %s\n", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		s.errors <- err
		return
	}
}

func (s *Server) Stop(wg *sync.WaitGroup) error {
	defer wg.Done()
	return s.Server.Shutdown(context.Background())
}

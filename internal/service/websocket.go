package service

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleGameReplay(w http.ResponseWriter, r *http.Request) {
	connID := uuid.New()
	connTime := time.Now()

	vars := mux.Vars(r)
	gameID, err := uuid.Parse(vars["game_id"])
	if err != nil {
		http.Error(w, "Missing Game ID: GET /replay/<game_id>", http.StatusNotFound)
		return
	}

	// db, err := database.Connection()
	// if err != nil {
	// 	log.Println(err.Error())
	// 	http.Error(w, "Error establishing database connection", http.StatusInternalServerError)
	//  return
	// }

	// game := db.Get

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Printf("[%s] New connection established", connID.String())

	ticker := time.NewTimer(5 * time.Second)

	stop := make(chan struct{})
	conn.SetCloseHandler(func(code int, text string) error {
		dur := time.Since(connTime)
		stop <- struct{}{}
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed")) // nolint
		log.Printf("[%s] Closed connection after %s", connID, dur)
		return nil
	})

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("[%s] Emitting message", connID)
				update := map[string]interface{}{
					"game_id":   gameID.String(),
					"timestamp": time.Now().Format(time.RFC3339Nano),
				}
				if err := conn.WriteJSON(update); err != nil {
					log.Println(err)
					conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())) // nolint
					return
				}
				ticker = time.NewTimer(1 * time.Second)
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil && !errors.Is(err, &websocket.CloseError{}) {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())) // nolint
			return
		}
	}
}

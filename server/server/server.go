package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/livingpool/top-down-shooter/game/game"
	"github.com/livingpool/top-down-shooter/game/pkg/player"
	"github.com/livingpool/top-down-shooter/game/util"
)

// GameServer maintains the set of active players
// and broadcasts game states at fixed intervals.
type GameServer struct {
	games    map[uuid.UUID]*game.Game // active games, each having one or more players
	serveMux *http.ServeMux           // serveMux routes endpoints to appropriate handlers
	mutex    *sync.Mutex
}

func NewGameServer() *GameServer {
	serveMux := http.NewServeMux()

	gs := &GameServer{
		games:    make(map[uuid.UUID]*game.Game),
		serveMux: serveMux,
		mutex:    &sync.Mutex{},
	}

	serveMux.HandleFunc("/", serveHome)
	serveMux.HandleFunc("/health", serveHealth)
	serveMux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		gs.create(w, r)
	})
	serveMux.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		gs.join(w, r)
	})
	// serveMux.HandleFunc("/start")

	return gs
}

func (gs *GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gs.serveMux.ServeHTTP(w, r)
}

// create creates a new game and a new player.
// create upgrades the connection to a websocket.
func (gs *GameServer) create(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query().Get("name")
	if playerName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no query param: name"))
		return
	}

	player := player.NewPlayer(playerName)
	game := game.NewGame(true)
	gs.games[game.ID] = game
	gs.subscribe(w, r, player, game.ID)

	log.Println("created new game:", game.ID, "playerName:", playerName)
}

// join adds the requested player to the game.
// join upgrades the connection to a websocket.
func (gs *GameServer) join(w http.ResponseWriter, r *http.Request) {
	log.Println("todo")
}

// subscribe accepts the websocket connetion and subcribes it to future game updates.
// It also listens for client updates and save them to a buffered channel.
func (gs *GameServer) subscribe(w http.ResponseWriter, r *http.Request, player *player.Player, gameId uuid.UUID) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("error upgrading conn to a websocket:", err)
		return
	}
	player.Conn = conn

	if err := gs.addSubscriber(player, gameId); err != nil {
		log.Println("error adding subscriber:", err)
		return
	}

	_, reader, err := conn.Reader(context.Background())
	if err != nil {
		log.Println("error calling conn.Reader:", err)
		return
	}

	for {
		data, err := io.ReadAll(reader)
		if err != nil {
			log.Println("connection closed:", err)
			return
		}

		var msg util.ClientUpdate
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Println("error unmarshaling data:", err)
			return
		}

		log.Println("accepted update:", msg)
		player.ClientUpdates = append(player.ClientUpdates, msg)
	}
}

// publish publishes each room's game state at fixed intervals to every subscriber in the room.
func (gs *GameServer) publish() {
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Yello!"))
}

func serveHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Game server is healthy!"))
}

func (gs *GameServer) addSubscriber(p *player.Player, gameId uuid.UUID) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	game, exists := gs.games[gameId]
	if !exists {
		return fmt.Errorf("room %v doesn't exist", gameId)
	}
	game.Players[p.ID] = p

	return nil
}

// deleteSubscriber deletes the subscriber from the room and closes its connection
func (gs *GameServer) deleteSubscriber(p *player.Player, gameId uuid.UUID) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	game, exists := gs.games[gameId]
	if !exists {
		return fmt.Errorf("room %v doesn't exist", gameId)
	}

	if _, exists := game.Players[p.ID]; !exists {
		return fmt.Errorf("room: %v doesn't contain player: %v", gameId, p.ID)
	}
	delete(game.Players, p.ID)

	err := p.Conn.Close(websocket.StatusNormalClosure, "delete request")
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}

	return nil
}

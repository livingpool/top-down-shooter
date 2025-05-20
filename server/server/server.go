package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
	logger   *slog.Logger
}

func NewGameServer() *GameServer {
	serveMux := http.NewServeMux()

	gs := &GameServer{
		games:    make(map[uuid.UUID]*game.Game),
		serveMux: serveMux,
		mutex:    &sync.Mutex{},
		logger:   slog.Default(),
	}

	serveMux.HandleFunc("/", serveHome)
	serveMux.HandleFunc("/health", serveHealth)
	serveMux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		gs.create(w, r)
	})
	serveMux.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		gs.join(w, r)
	})
	// serveMux.HandleFunc("/delete")
	// serveMux.HandleFunc("/start")

	return gs
}

func (gs *GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gs.serveMux.ServeHTTP(w, r)
}

// create creates a new game and a new player and returns the associated ids.
// The client should call /join after /create.
func (gs *GameServer) create(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query().Get("name")
	if playerName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no query param: name"))
		return
	}

	player := player.NewPlayer(playerName)
	game := game.NewGame(true)
	if err := gs.addPlayer(player, game); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(util.CreatePlayerResp{
		PlayerId: player.ID.String(),
		GameId:   game.ID.String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

	gs.logger.Info("new game created", "# of games", len(gs.games), "player name", playerName)
}

// join adds the requested player to the game.
// join upgrades the connection to a websocket.
func (gs *GameServer) join(w http.ResponseWriter, r *http.Request) {
	playerId := r.URL.Query().Get("player_id")
	gameId := r.URL.Query().Get("game_id")

	game, player, err := gs.getPlayer(playerId, gameId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("player id or game id is not uuid"))
	}

	gs.subscribe(w, r, game, player)
}

// subscribe accepts the websocket connetion and subcribes it to future game updates.
// It also listens for client updates and save them to a buffered channel.
func (gs *GameServer) subscribe(w http.ResponseWriter, r *http.Request, game *game.Game, player *player.Player) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		gs.logger.Error("error upgrading conn to a websocket: %v", "err", err)
		return
	}
	player.Conn = conn

	_, reader, err := conn.Reader(context.Background())
	if err != nil {
		gs.logger.Error("error creating websocket reader", "err", err)
		return
	}

	for {
		data, err := io.ReadAll(reader)
		if err != nil {
			gs.logger.Debug("connection closed", "err", err)
			return
		}

		var msg util.ClientUpdate
		if err := json.Unmarshal(data, &msg); err != nil {
			gs.logger.Error("error unmarshaling client data", "err", err)
			return
		}

		slog.Debug("accepted client update", "msg", msg)
		player.ClientUpdates = append(player.ClientUpdates, msg)

		// poc
		if err := gs.processInput(game); err != nil {
			gs.logger.Error("error processing input", "err", err)
		}
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

func (gs *GameServer) getPlayer(playerIdStr, gameIdStr string) (*game.Game, *player.Player, error) {
	playerId, err := uuid.Parse(playerIdStr)
	if err != nil {
		return nil, nil, fmt.Errorf("playerId is not uuid")
	}
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		return nil, nil, fmt.Errorf("gameId is not uuid")
	}

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	game, exists := gs.games[gameId]
	if !exists {
		return nil, nil, fmt.Errorf("game %v does not exist", gameId)
	}

	player, exists := game.Players[playerId]
	if !exists {
		return nil, nil, fmt.Errorf("player %v does not exist", playerId)
	}

	return game, player, nil
}

func (gs *GameServer) addPlayer(player *player.Player, game *game.Game) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if _, exists := gs.games[game.ID]; exists {
		return fmt.Errorf("game %v exists", game.ID)
	} else {
		gs.games[game.ID] = game
	}

	if _, exists := game.Players[player.ID]; exists {
		return fmt.Errorf("player %v exists", player.ID)
	} else {
		game.Players[player.ID] = player
	}

	return nil
}

// deleteSubscriber deletes the subscriber from the room and closes its connection
func (gs *GameServer) deleteSubscriber(playerIdStr, gameIdStr string) error {
	playerId, err := uuid.Parse(playerIdStr)
	if err != nil {
		return fmt.Errorf("playerId is not uuid")
	}
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		return fmt.Errorf("gameId is not uuid")
	}

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	game, exists := gs.games[gameId]
	if !exists {
		return fmt.Errorf("game %v doesn't exist", gameId)
	}

	player, exists := game.Players[playerId]
	if !exists {
		return fmt.Errorf("game: %v doesn't contain player: %v", gameId, playerId)
	}
	delete(game.Players, playerId)

	err = player.Conn.Close(websocket.StatusNormalClosure, "delete request")
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}

	return nil
}

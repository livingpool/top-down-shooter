package server

import (
	"net/http"
	"sync"

	"github.com/livingpool/top-down-shooter/server/model"
)

// GameServer maintains the set of active players and rooms
// and broadcasts game states at fixed intervals.
type GameServer struct {
	subscribers map[string]*model.Subscriber // registered clients
	rooms       map[string]*model.Room       // created rooms
	register    chan *model.Subscriber       // register requests from clients
	unregister  chan *model.Subscriber       // unregister requests from clients
	serveMux    *http.ServeMux               // serveMux routes endpoints to appropriate handlers
	mutex       *sync.RWMutex
}

func NewGameServer() *GameServer {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", serveHome)
	serveMux.HandleFunc("/health", serveHealth)

	return &GameServer{
		subscribers: make(map[string]*model.Subscriber),
		rooms:       make(map[string]*model.Room),
		register:    make(chan *model.Subscriber),
		unregister:  make(chan *model.Subscriber),
		serveMux:    serveMux,
		mutex:       &sync.RWMutex{},
	}
}

func (gs *GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gs.serveMux.ServeHTTP(w, r)
}

// subscribe accepts the websocket connetion and subcribes it to future game updates.
// It listens for client updates and save them to a buffered channel.
func (gs *GameServer) subscribe() {
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

func (gs *GameServer) addSubscriber(s *model.Subscriber, roomId string) {
	room := gs.rooms[roomId]
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	room.Subscribers[s.Player.ID] = s
}

func (gs *GameServer) deleteSubscriber(s *model.Subscriber, roomId string) {
	room := gs.rooms[roomId]
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	delete(room.Subscribers, s.Player.ID)
}

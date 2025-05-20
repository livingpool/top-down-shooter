package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coder/websocket"
	"github.com/livingpool/top-down-shooter/game/util"
)

// TODO: im not entirely sure of the behavior if i set up the game server and then run all tests in parallel
// TODO: tj's snapshot testing sounds cool; can i do it for testing games states?
func TestGameServer(t *testing.T) {
}

func setupTest() (url string, closeFn func()) {
	gs := NewGameServer()
	server := httptest.NewServer(gs)
	return server.URL, server.Close
}

// websocket client for testing
type client struct {
	t        *testing.T
	url      string
	conn     *websocket.Conn
	gameId   string
	playerId string
}

func newClient(t *testing.T, ctx context.Context, url string, playerName string) *client {
	t.Helper()

	resp, err := http.Get(url + "/create")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var createResp util.CreatePlayerResp
	if err := json.Unmarshal(data, &createResp); err != nil {
		t.Fatal(err)
	}

	conn, _, err := websocket.Dial(ctx, url+"/join", nil)
	if err != nil {
		t.Fatal(err)
	}

	return &client{t: t, url: url, conn: conn, gameId: createResp.GameId, playerId: createResp.PlayerId}
}

func (cl *client) publishMsg(ctx context.Context, msg []byte) error {
	cl.t.Helper()
	return cl.conn.Write(ctx, websocket.MessageText, msg)
}

func (cl *client) nextMessage(ctx context.Context) (string, error) {
	cl.t.Helper()

	typ, data, err := cl.conn.Read(ctx)
	if err != nil {
		return "", err
	}

	if typ != websocket.MessageText {
		cl.conn.Close(websocket.StatusUnsupportedData, "expected text message")
		return "", fmt.Errorf("expected text message but got %v", typ)
	}

	return string(data), nil
}

func (cl *client) close() error {
	cl.t.Helper()
	return cl.conn.Close(websocket.StatusNormalClosure, "")
}

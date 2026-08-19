package services

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type GameRoom struct {
	PlayerX *websocket.Conn
	PlayerO *websocket.Conn
}

type WebsocketHub struct {
	mu    sync.RWMutex
	games map[uuid.UUID]*GameRoom
}

func NewWebsocketHub() *WebsocketHub {
	return &WebsocketHub{
		games: make(map[uuid.UUID]*GameRoom),
	}
}

func (h *WebsocketHub) AddPlayer(
	gameID uuid.UUID,
	playerID uuid.UUID,
	userX uuid.UUID,
	userO uuid.UUID,
	conn *websocket.Conn,
) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.games[gameID]
	if !exists {
		room = &GameRoom{}
		h.games[gameID] = room
	}

	if playerID == userX {
		room.PlayerX = conn
		return
	}

	if playerID == userO {
		room.PlayerO = conn
	}
}

func (h *WebsocketHub) GetRoom(
	gameID uuid.UUID,
) *GameRoom {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.games[gameID]
}

func (h *WebsocketHub) RemovePlayer(
	gameID uuid.UUID,
	playerID uuid.UUID,
	userX uuid.UUID,
	userO uuid.UUID,
) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.games[gameID]
	if !exists {
		return
	}

	if playerID == userX {
		room.PlayerX = nil
	}

	if playerID == userO {
		room.PlayerO = nil
	}

	if room.PlayerX == nil && room.PlayerO == nil {
		delete(h.games, gameID)
	}
}

package services

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type GameRoom struct {
	PlayerX *websocket.Conn
	PlayerO *websocket.Conn
	writeMu sync.Mutex
}

func (h *WebsocketHub) SendToConnection(
	gameID uuid.UUID,
	conn *websocket.Conn,
	payload any,
) error {
	h.mu.RLock()
	room := h.games[gameID]
	if room == nil {
		h.mu.RUnlock()
		return nil
	}
	h.mu.RUnlock()

	if conn == nil {
		return nil
	}

	room.writeMu.Lock()
	defer room.writeMu.Unlock()
	return conn.WriteJSON(payload)
}

func (h *WebsocketHub) Broadcast(gameID uuid.UUID, payload any) []error {
	h.mu.RLock()
	room := h.games[gameID]
	if room == nil {
		h.mu.RUnlock()
		return nil
	}
	playerX := room.PlayerX
	playerO := room.PlayerO
	h.mu.RUnlock()

	room.writeMu.Lock()
	defer room.writeMu.Unlock()

	var errs []error
	if playerX != nil {
		if err := playerX.WriteJSON(payload); err != nil {
			errs = append(errs, err)
		}
	}
	if playerO != nil && playerO != playerX {
		if err := playerO.WriteJSON(payload); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
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

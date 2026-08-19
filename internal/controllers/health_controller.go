package controllers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HealthController struct {
	client    *mongo.Client
	startedAt time.Time
}

type healthResponse struct {
	Status   string         `json:"status"`
	Uptime   string         `json:"uptime"`
	Database databaseHealth `json:"database"`
}

type databaseHealth struct {
	Status    string `json:"status"`
	LatencyMS int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

func NewHealthController(client *mongo.Client, startedAt time.Time) *HealthController {
	return &HealthController{client: client, startedAt: startedAt}
}

func (c *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	started := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	err := c.client.Ping(ctx, nil)

	response := healthResponse{
		Status: "ok",
		Uptime: time.Since(c.startedAt).Round(time.Second).String(),
		Database: databaseHealth{
			Status:    "ok",
			LatencyMS: time.Since(started).Milliseconds(),
		},
	}
	status := http.StatusOK
	if err != nil {
		status = http.StatusServiceUnavailable
		response.Status = "degraded"
		response.Database.Status = "error"
		response.Database.Error = err.Error()
		slog.Error("database health check failed", "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		slog.Error("failed to encode health response", "error", encodeErr)
	}
}

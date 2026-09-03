package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AlcantuMurilo77/tic/internal/controllers"
	"github.com/AlcantuMurilo77/tic/internal/database"
	"github.com/AlcantuMurilo77/tic/internal/middleware"
	"github.com/AlcantuMurilo77/tic/internal/repository"
	"github.com/AlcantuMurilo77/tic/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	startedAt := time.Now()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC)

	if err := godotenv.Load(); err != nil {
		slog.Info(".env file not loaded; using environment variables", "error", err)
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		slog.Error("MONGODB_URI is required")
		os.Exit(1)
	}

	client, err := database.Connect(uri)
	if err != nil {
		slog.Error("failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to MongoDB")

	db := client.Database(os.Getenv("DATABASE_NAME"))
	gameRepository := repository.NewGameRepository(db)
	gameService := services.NewGameService(gameRepository)
	gameController := controllers.NewGameController(gameService)

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	webSocketService := services.NewWebSocketService()
	webSocketHub := services.NewWebsocketHub()
	webSocketController := controllers.NewWebSocketController(webSocketService, gameService, webSocketHub)
	healthController := controllers.NewHealthController(client, startedAt)

	rematchRepository := repository.NewRematchRepository(db)
	if err := rematchRepository.EnsureIndexes(context.Background()); err != nil {
		slog.Error("failed to create rematch indexes", "error", err)
		os.Exit(1)
	}
	rematchService := services.NewRematchService(rematchRepository, gameRepository)
	rematchController := controllers.NewRematchController(rematchService, webSocketHub)

	mux := http.NewServeMux()
	mux.HandleFunc("/games/join", gameController.Join)
	mux.HandleFunc("/games", gameController.Create)
	mux.HandleFunc("/games/get_all", gameController.FindAll)
	mux.HandleFunc("/users", userController.Create)
	mux.HandleFunc("/users/get_all", userController.FindAll)
	mux.HandleFunc("/ws", webSocketController.Connect)
	mux.HandleFunc("/healthz", healthController.Check)
	mux.HandleFunc("/games/rematch", rematchController.Request)
	mux.HandleFunc("/games/rematch/accept", rematchController.Accept)

	// allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	// if allowedOrigins == "" {
	allowedOrigins := "*"
	// }
	handler := middleware.Logging(
		middleware.Recoverer(
			middleware.Timeout(10*time.Second,
				middleware.CORS(allowedOrigins, mux),
			),
		),
	)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info("server started", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server stopped unexpectedly", "error", err)
		os.Exit(1)
	}

}

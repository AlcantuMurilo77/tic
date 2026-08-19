package main

import (
	"context"
	"log"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}
	uri := os.Getenv("MONGODB_URI")

	client, err := database.Connect(uri)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(os.Getenv("DATABASE_NAME"))
	gameRepository := repository.NewGameRepository(db) //gameRepository
	gameService := services.NewGameService(gameRepository)
	gameController := controllers.NewGameController(gameService)

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	webSocketService := services.NewWebSocketService()
	webSocketHub := services.NewWebsocketHub()
	webSocketController := controllers.NewWebSocketController(webSocketService, gameService, webSocketHub)

	mux := http.NewServeMux()
	mux.HandleFunc("/games/join", gameController.Join)
	mux.HandleFunc("/games", gameController.Create)
	mux.HandleFunc("/games/get_all", gameController.FindAll)
	mux.HandleFunc("/users", userController.Create)
	mux.HandleFunc("/users/get_all", userController.FindAll)
	mux.HandleFunc("/ws", webSocketController.Connect)

	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173"
	}
	handler := middleware.CORS(allowedOrigins, mux)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}

}

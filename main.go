package main

import(
	"context"
	"net/http"
	"log"
	"time"
	"os"
	"github.com/AlcantuMurilo77/tic/internal/database"
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

	defer func(){
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

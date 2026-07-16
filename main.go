package main

import(
	"net/http"
	"log"
)


func main() {
	log.Println("Server running on :8080")
	game := NewTicTacToe(3)

	game.Move(0,0, 1)
	game.Move(1,0,1)
	game.Move(2, 0, 1)
	result := game.Move(2, 0, 2)


	log.Println(game.Rows)
	log.Println(game.Cols)
	log.Println(result)
	http.ListenAndServe(":8080", nil)
}

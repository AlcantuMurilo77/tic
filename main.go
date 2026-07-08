package main

import(
	"net/http"
	"log"
)


func main() {
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/Go_CleanArch/infrastructure/db"
	"github.com/Go_CleanArch/infrastructure/server"
)

func main() {

	fmt.Println("Starting Server...")
	db.Init()
	server.Init()
	http.ListenAndServe(":8080", nil)
}

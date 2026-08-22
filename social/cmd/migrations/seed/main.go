package main

import (
	"log"
	"main/internal/db"
	"main/internal/env"
	"main/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgresql://user:Nathansno1@localhost/socialnetwork?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	stores := store.NewStorage(conn)

	db.Seed(stores, conn)
}

package main

import (
	"log"

	"github.com/ccichielo/gobank/pkg"
)

func main() {
	store, err := pkg.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := pkg.NewAPIServer(":3000", store)
	server.Run()
}

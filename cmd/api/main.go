package main

import (
	"log"

	"github.com/Nasaee/go-social/internal/env"
	"github.com/Nasaee/go-social/internal/store"
)

func main() {
	env.Init()

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	// log.Fatal() is a shortcut for os.Exit(1) จอ error ที่ต้องหยุดระบบ
	log.Fatal(app.run(mux))
}

package main

import (
	"log"

	"github.com/Nasaee/go-social/internal/db"
	"github.com/Nasaee/go-social/internal/env"
	"github.com/Nasaee/go-social/internal/store"
)

func main() {
	env.Init()

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://postgres@localhost/go_social?sslmode=disable"), // put the real connection string here it will pick DB_ADDR from .env first and if it's not set then fallback to second
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		// log.Fatal() เหมาะกับ “critical startup error” เช่น DB ต่อไม่ได้, พอร์ตซ้ำ, env ผิด เพราะเรา ไม่ต้องการให้รันต่อ แล้วจะให้ระบบ exit ไปเลย
		log.Fatal(err)  

		// log.Panic() เหมาะกับ “unexpected runtime error” ที่ยังอยากเห็น stack trace หรือจะจับได้ด้วย recove()
		// log.Panic(err)
	}

	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	// log.Fatal() is a shortcut for os.Exit(1) จอ error ที่ต้องหยุดระบบ
	log.Fatal(app.run(mux))
}

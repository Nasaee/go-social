package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Nasaee/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
}

func (app *application) mount() *chi.Mux {
	// ServeMux is an HTTP request multiplexer.
	r := chi.NewRouter()

	// ref: https://github.com/go-chi/chi
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})
	return r
}

func (app *application) run(mux *chi.Mux) error {
	/*
		IdleTimeout = เวลาสูงสุดที่ server จะ รอให้ connection เดิมว่างเปล่า (idle) ก่อนจะปิด connection นั้นลง

		พูดง่าย ๆ คือ
		ถ้า client เปิด connection ไว้ แต่ “ไม่ได้ส่ง request ใหม่มาอีกเลย” หลังจากเวลาที่กำหนดใน IdleTimeout → server จะบอกว่า “โอเค ขอตัดสายละนะ 👋”
	*/
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30, // 30 seconds
		ReadTimeout:  time.Second * 10, // 30 seconds
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}

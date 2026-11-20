package main

import (
	"log"
	"net/http"
	"time"

	"github.com/maillunex-lang/404Eco/backend/internal/api"
	"github.com/maillunex-lang/404Eco/backend/internal/config"
	"github.com/maillunex-lang/404Eco/backend/internal/economy"
)

func main() {
	cfg := config.Load()
	engine := economy.NewEngine()
	srv := api.New(engine)

	go func() {
		t := time.NewTicker(cfg.TickInterval)
		defer t.Stop()
		for range t.C {
			snap := engine.Tick()
			srv.BroadcastSnapshot(snap)
		}
	}()

	handler := srv.Handler(cfg.CORSOrigin)
	log.Printf("Bashkortostan Eco API on %s (tick=%s, memory=%v)", cfg.HTTPAddr, cfg.TickInterval, cfg.UseMemory)
	if err := http.ListenAndServe(cfg.HTTPAddr, handler); err != nil {
		log.Fatal(err)
	}
}

package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/maillunex-lang/404Eco/backend/internal/economy"
	"github.com/maillunex-lang/404Eco/backend/internal/models"
	"github.com/rs/cors"
)

type Server struct {
	engine *economy.Engine
	mux    *http.ServeMux

	hubMu sync.Mutex
	hub   map[chan []byte]struct{}
}

func New(engine *economy.Engine) *Server {
	s := &Server{
		engine: engine,
		mux:    http.NewServeMux(),
		hub:    map[chan []byte]struct{}{},
	}
	s.routes()
	return s
}

func (s *Server) Handler(corsOrigin string) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{corsOrigin},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	if corsOrigin == "*" {
		c = cors.AllowAll()
	}
	return c.Handler(s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("GET /api/snapshot", s.snapshot)
	s.mux.HandleFunc("GET /api/industries", s.industries)
	s.mux.HandleFunc("POST /api/register", s.register)
	s.mux.HandleFunc("POST /api/firms", s.createFirm)
	s.mux.HandleFunc("POST /api/firms/action", s.firmAction)
	s.mux.HandleFunc("POST /api/orders", s.placeOrder)
	s.mux.HandleFunc("GET /api/stream", s.stream)
	s.mux.HandleFunc("POST /api/nakama/sync", s.nakamaSync)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"ok": true, "service": "bashkir-eco"})
}

func (s *Server) snapshot(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, s.engine.Snapshot())
}

func (s *Server) industries(w http.ResponseWriter, _ *http.Request) {
	type item struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	}
	out := make([]item, 0, len(economy.IndustryLabels))
	for id, label := range economy.IndustryLabels {
		out = append(out, item{ID: id, Label: label})
	}
	writeJSON(w, out)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NakamaID string `json:"nakamaId"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if body.Name == "" {
		body.Name = "Предприниматель"
	}
	p, err := s.engine.RegisterPlayer(body.NakamaID, body.Name)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, p)
}

func (s *Server) createFirm(w http.ResponseWriter, r *http.Request) {
	var req economy.CreateFirmReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	firm, err := s.engine.CreateFirm(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	s.broadcast()
	writeJSON(w, firm)
}

func (s *Server) firmAction(w http.ResponseWriter, r *http.Request) {
	var req economy.FirmActionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	firm, err := s.engine.FirmAction(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	s.broadcast()
	writeJSON(w, firm)
}

func (s *Server) placeOrder(w http.ResponseWriter, r *http.Request) {
	var req economy.OrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	order, err := s.engine.PlaceOrder(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	s.broadcast()
	writeJSON(w, order)
}

func (s *Server) nakamaSync(w http.ResponseWriter, r *http.Request) {
	// Lightweight bridge for Nakama RPCs — returns current tick + player cash.
	var body struct {
		NakamaID string `json:"nakamaId"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	snap := s.engine.Snapshot()
	writeJSON(w, map[string]any{
		"tick":  snap.World.Tick,
		"world": snap.World,
		"ok":    true,
	})
}

func (s *Server) stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", 500)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan []byte, 4)
	s.hubMu.Lock()
	s.hub[ch] = struct{}{}
	s.hubMu.Unlock()
	defer func() {
		s.hubMu.Lock()
		delete(s.hub, ch)
		s.hubMu.Unlock()
		close(ch)
	}()

	// initial
	init, _ := json.Marshal(s.engine.Snapshot())
	_, _ = w.Write([]byte("event: snapshot\ndata: "))
	_, _ = w.Write(init)
	_, _ = w.Write([]byte("\n\n"))
	flusher.Flush()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-ch:
			_, _ = w.Write([]byte("event: snapshot\ndata: "))
			_, _ = w.Write(msg)
			_, _ = w.Write([]byte("\n\n"))
			flusher.Flush()
		case <-ticker.C:
			_, _ = w.Write([]byte(": ping\n\n"))
			flusher.Flush()
		}
	}
}

func (s *Server) BroadcastSnapshot(snap models.Snapshot) {
	data, err := json.Marshal(snap)
	if err != nil {
		return
	}
	s.hubMu.Lock()
	defer s.hubMu.Unlock()
	for ch := range s.hub {
		select {
		case ch <- data:
		default:
		}
	}
}

func (s *Server) broadcast() {
	s.BroadcastSnapshot(s.engine.Snapshot())
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

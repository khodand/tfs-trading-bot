package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/internal/services"
	"tfs-trading-bot/internal/services/algorithms"
)

type Server struct {
	service services.TradingService
}

func NewServer(trader services.TradingService) *Server {
	return &Server{
		service: trader,
	}
}

func (s *Server) Router() chi.Router {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Post("/trade/{tickerSymbol}", s.TradeTicker)
	root.Post("/algo/{algorithm}/{period}", s.ChangeAlgo)

	return root
}

func (s *Server) TradeTicker(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "tickerSymbol")
	s.service.TradeTicker(domain.TickerSymbol(symbol))
}

func (s *Server) ChangeAlgo(w http.ResponseWriter, r *http.Request) {
	algoName := chi.URLParam(r, "algorithm")
	period := chi.URLParam(r, "period")
	atoi, err := strconv.Atoi(period)
	if err != nil {
		writeError(w, "Period must be int.")
	}
	var algo services.TradingAlgorithm
	switch algoName {
	case "EMA":
		algo = algorithms.NewEMAAlgo(atoi)
	default:
		writeError(w, "No such algo.")
	}
	s.service.ChangeAlgo(algo)
}

func writeError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(msg))
}

package server

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type RatesClient interface {
	GetLatestRatesForUSD() (map[string]float64, error)
}

type Server struct {
	router         *gin.Engine
	exchangeClient RatesClient
}

func NewServer(client RatesClient) *Server {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &Server{
		router:         router,
		exchangeClient: client,
	}

	server.setupRoutes()

	return server
}

func (s *Server) Start() error {
	return s.router.Run(":8080")
}

func (s *Server) setupRoutes() {
	s.router.GET("/rates", s.getRates)
	s.router.GET("/exchange", s.getCryptoExchange)
}

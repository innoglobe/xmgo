package app

import (
	"context"
	"errors"
	"fmt"
	eventservice "github.com/innoglobe/xmgo/internal/service"
	"github.com/innoglobe/xmgo/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/innoglobe/xmgo/internal/config"
	"github.com/innoglobe/xmgo/pkg/logger"
)

type App interface {
	Run() error
}

type app struct {
	//Router         server.RouterInterface
	Server         *http.Server
	Config         *config.Config
	CompanyUseCase usecase.CompanyUsecaseInterface
	Logger         logger.LoggerInterface
	KafkaProducer  *eventservice.KafkaProducer
}

func NewApp(cfg *config.Config, companyUsecase usecase.CompanyUsecaseInterface, log logger.LoggerInterface, router *gin.Engine, kafkaProducer *eventservice.KafkaProducer) (App, error) {
	return &app{
		//Router:         router,
		Server:         &http.Server{Addr: fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port), Handler: router},
		CompanyUseCase: companyUsecase,
		Config:         cfg,
		Logger:         log,
		KafkaProducer:  kafkaProducer,
	}, nil
}

func (a *app) Run() error {
	// Start server in a goroutine
	go func() {
		var err error
		if a.Config.Server.SSL.Enabled {
			err = a.Server.ListenAndServeTLS(a.Config.Server.SSL.CertFile, a.Config.Server.SSL.KeyFile)
		} else {
			err = a.Server.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Error(fmt.Sprintf("Couldn't start server: %v", err))
		}
	}()

	// Wait for shutdown signal
	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q
	a.Logger.Info("Shutdown signal received, shutting down server...")

	// Close the Kafka producer
	if err := a.KafkaProducer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka producer: %v", err)
	} else {
		a.Logger.Info("Kafka producer closed")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Config.Server.Timeout)*time.Second)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	} else {
		a.Logger.Info("Server gracefully shutdown")
	}

	a.Logger.Info("Server exiting")
	return nil
}

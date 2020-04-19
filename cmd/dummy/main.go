package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chamilto/dummy/internal/config"
	"github.com/chamilto/dummy/internal/db"
	"github.com/chamilto/dummy/internal/handlers"
	"github.com/sirupsen/logrus"
)

func init() {
	config.LoadEnv()
}

func main() {
	c := config.NewConfig()
	db := db.NewDB(c)
	ctlr := handlers.NewHandlerController(c, db)
	r := handlers.NewRouter(ctlr)

	server := &http.Server{
		Handler:      r,
		Addr:         c.Bind,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Error starting the server: %s\n", err)
		}
	}()

	logrus.Printf("Server started on %s", server.Addr)

	<-done
	logrus.Print("Server stopped")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		ctlr.DB.Close()
		cancel()
	}()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.Fatal(err)
	}

	logrus.Print("Server exited")

}

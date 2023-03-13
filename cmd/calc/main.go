package main

import (
	"github.com/go-chi/httplog"
	"main/config"
	"main/internal/handlers"
	"main/internal/logger"
	"main/internal/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)
	cfg := config.New()

	log := logger.New(httplog.NewLogger("market", httplog.Options{
		Concise: true,
	}))

	h := handlers.New(usecase.New(log), log)

	http.HandleFunc("/", h.MainHandler)
	http.HandleFunc("/res", h.ResultHandler)
	http.Handle("/static", http.FileServer(http.Dir("static/")))
	go http.ListenAndServe(cfg.Host, nil)

	log.Info("calc server up on " + cfg.Host)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutdown market ...")
}

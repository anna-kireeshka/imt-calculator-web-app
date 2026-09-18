package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/go-telegram/bot"
)

type Application struct {
	botToken string
	pages    http.Handler
	b        *bot.Bot
}

func NewTg(botToken string, pages http.Handler) (*Application, error) {
	var app = &Application{
		botToken: botToken,
		pages:    pages,
	}

	var b, errBot = bot.New(botToken)
	if errBot != nil {
		return nil, fmt.Errorf("failed to create bot: %w", errBot)
	}

	app.b = b

	return app, nil
}

func (app *Application) Run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()

	var mux = http.NewServeMux()
	mux.HandleFunc("/api/open", app.handlerAPIOpen)
	mux.Handle("/", app.pages)

	var server = http.Server{
		Handler: mux,
	}

	wg.Add(1)
	go func() {
		log.Printf("bot starting")
		app.b.Start(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()
		log.Printf("server stopping")
		server.Shutdown(context.Background())
		wg.Done()
	}()

	log.Printf("server started at %s", ln.Addr().String())
	var errServe = server.Serve(ln)
	if errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
		log.Printf("error serve: %v", errServe)
	}
	cancel()
}

func (app *Application) handlerAPIOpen(rw http.ResponseWriter, req *http.Request) {
	user, ok := bot.ValidateWebappRequest(req.URL.Query(), app.botToken)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("%v", user)
}

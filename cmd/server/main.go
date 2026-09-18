package main

import (
	"app/imt-calculator-web-app/internal/config"
	"app/imt-calculator-web-app/internal/database"
	"app/imt-calculator-web-app/internal/server"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

func main() {
	var cfg, errConfig = config.LoadConfig()

	if errConfig != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", errConfig)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var repo, err = database.NewRepository(ctx, cfg.DataBaseUrl)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}

	defer repo.Close(context.Background())

	fmt.Println("Сервер запущен")

	srv := server.New(cfg, repo)

	if err := runTg(ctx, cancel, cfg, srv.Router()); err != nil {
		log.Printf("error run: %v", err)
		os.Exit(1)
	}
}

func runTg(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, pages http.Handler) error {
	ln, errLn := net.Listen("tcp", cfg.ListenAddress)
	if errLn != nil {
		return fmt.Errorf("failed to listen: %w", errLn)
	}
	defer ln.Close()

	app, errApp := server.NewTg(cfg.TgBotToken, pages)
	if errApp != nil {
		return fmt.Errorf("failed to create application: %w", errApp)
	}

	log.Printf("start")

	var wg sync.WaitGroup

	wg.Add(1)
	go app.Run(ctx, cancel, &wg, ln)

	<-ctx.Done()

	wg.Wait()

	log.Printf("done")

	return nil
}

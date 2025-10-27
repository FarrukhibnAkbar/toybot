package app

import (
	"context"
	"fmt"

	"toybot/internal/bot"
	dbq "toybot/internal/db"
	"toybot/internal/repository"

	"go.uber.org/zap"
)

type App struct {
	cfg    *Config
	logger *zap.Logger
	db     *repository.Postgres
	bot    *bot.Bot
}

func New() (*App, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config load error: %w", err)
	}

	logger, _ := zap.NewProduction()

	db, err := repository.NewPostgres(context.Background(), cfg.DatabaseURL, logger)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	queries := dbq.New(db.Pool)

	b, err := bot.NewBot(cfg.BotToken, cfg.AllowedUsers, logger, queries)
	if err != nil {
		return nil, fmt.Errorf("bot init error: %w", err)
	}

	app := &App{
		cfg:    cfg,
		logger: logger,
		db:     db,
		bot:    b,
	}

	return app, nil
}

func (a *App) Run() {
	a.logger.Info("🚀 Toyshop Bot is running...")
	a.bot.Start()
}

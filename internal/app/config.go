package app

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	DatabaseURL  string
	AllowedUsers []int64
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // quietly load .env if exists

	cfg := &Config{}
	cfg.BotToken = os.Getenv("BOT_TOKEN")
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")

	usersRaw := os.Getenv("ALLOWED_USERS")
	for _, v := range strings.Split(usersRaw, ",") {
		if v == "" {
			continue
		}
		id, _ := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		cfg.AllowedUsers = append(cfg.AllowedUsers, id)
	}

	return cfg, nil
}

package bot

import (
	"log"
	"toybot/internal/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	api          *tgbotapi.BotAPI
	allowedUsers []int64
	logger       *zap.Logger
	repo         *db.Queries
	sessions     map[int64]*Session
}

func NewBot(token string, allowedUsers []int64, logger *zap.Logger, repo *db.Queries) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = false
	logger.Info("✅ Connected to Telegram Bot API", zap.String("bot_username", botAPI.Self.UserName))

	return &Bot{
		api:          botAPI,
		allowedUsers: allowedUsers,
		logger:       logger,
		repo:         repo,
		sessions:     make(map[int64]*Session),
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		if !b.isAllowed(userID) {
			b.reply(update.Message.Chat.ID, "❌ Sizga bu botdan foydalanishga ruxsat yo‘q.")
			continue
		}

		text := update.Message.Text

		if session, ok := b.sessions[userID]; ok && session.State != StateNone {
			b.handleKirim(update.Message)
			continue
		}

		switch text {
		case "/start":
			b.reply(update.Message.Chat.ID, "👋 Assalomu alaykum!\nBu ToyShop ombor botidir.")
		case "/help":
			b.reply(update.Message.Chat.ID, "📋 Mavjud komandalar:\n/start - salomlashish\n/help - yordam")
		case "/kirim":
			b.handleKirim(update.Message)
		default:
			b.reply(update.Message.Chat.ID, "⚠️ Noma’lum komanda. /help buyrug‘ini yozing.")
		}
	}
}

func (b *Bot) isAllowed(id int64) bool {
	for _, u := range b.allowedUsers {
		if u == id {
			return true
		}
	}
	return false
}

func (b *Bot) reply(chatID int64, msg string) {
	message := tgbotapi.NewMessage(chatID, msg)
	if _, err := b.api.Send(message); err != nil {
		log.Printf("send message error: %v", err)
	}
}

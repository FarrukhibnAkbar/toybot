package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleStep(msg *tgbotapi.Message) {
	session := b.getSession(msg.From.ID)

	log.Printf("User %d, Current State: %v, Text: %s",
		msg.From.ID, session.State, msg.Text)

	switch session.State {
	case StateKirimName, StateKirimQty, StateKirimBuy, StateKirimSell, StateKirimExistsChoice:
		b.handleKirim(msg)
	case StateSellName, StateSellQty:
		b.processSotishStep(msg)
	default:
		b.reply(msg.Chat.ID, "⚠️ Iltimos, komandani /kirim yoki /sotish bilan boshlang.")
	}
}

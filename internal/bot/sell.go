package bot

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleSotish(msg *tgbotapi.Message) {
	userID := msg.From.ID
	session := b.getSession(userID)
	session.State = StateSellName
	b.reply(msg.Chat.ID, "🛒 Sotish jarayoni boshlandi.\n\nIltimos, tovar nomini kiriting:")
}

func (b *Bot) processSotishStep(msg *tgbotapi.Message) {
	userID := msg.From.ID
	session := b.getSession(userID)
	text := msg.Text

	users, err := b.repo.GetUserDataByTgID(context.Background(), int64(userID))
	if err != nil {
		b.reply(msg.Chat.ID, "❌ Foydalanuvchi ma’lumotlarini olishda xatolik.")
		return
	}

	switch session.State {
	case StateSellName:
		session.TempName = text
		session.State = StateSellQty
		b.reply(msg.Chat.ID, "Necha dona sotmoqchisiz?")

	case StateSellQty:
		qty, err := strconv.ParseFloat(text, 64)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Noto‘g‘ri son kiritildi. Qayta urinib ko‘ring:")
			return
		}

		session.TempQty = qty

		// Ombordagi mavjud tovarni olish
		item, err := b.repo.GetItemByName(context.Background(), session.TempName)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Bunday tovar topilmadi.")
			return
		}

		if item.Quantity < qty {
			b.reply(msg.Chat.ID, fmt.Sprintf("⚠️ Omborda faqat %.2f dona bor.", item.Quantity))
			session.State = StateSellQty
			return
		}

		// Yangi miqdorni yangilash
		newQty := item.Quantity - qty
		err = b.repo.UpdateItemQuantity(context.Background(), int(item.ID), newQty)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Ma’lumotni yangilashda xatolik.")
			return
		}

		// Sotuv tarixiga yozish
		err = b.repo.CreateSale(context.Background(), int(item.ID), qty, item.SellPrice, int(users.ID))
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Sotuvni saqlashda xatolik.")
			return
		}
		delete(b.sessions, userID)

		b.reply(msg.Chat.ID, fmt.Sprintf("✅ %.2f dona '%s' sotildi.\nOmborda %.2f dona qoldi.", qty, item.Name, newQty))
	}
}

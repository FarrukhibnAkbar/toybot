package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleKirim — foydalanuvchi bilan interaktiv tovar kiritish jarayoni
func (b *Bot) handleKirim(msg *tgbotapi.Message) {
	userID := msg.From.ID
	session, ok := b.sessions[userID]

	// 🔍 Debug log - session mavjudligini tekshirib log qilish
	if ok {
		log.Printf("User %d, Session exists: true, State: %v, Text: %s",
			userID, session.State, msg.Text)
	} else {
		log.Printf("User %d, Session exists: false, Text: %s",
			userID, msg.Text)
	}

	// Agar session mavjud bo'lmasa — yangi kirim jarayonini boshlaymiz
	if !ok {
		b.sessions[userID] = &Session{State: StateKirimName}
		b.reply(msg.Chat.ID, "📦 Tovar nomini kiriting:")
		return
	}

	text := strings.TrimSpace(msg.Text)

	switch session.State {
	case StateKirimName:
		log.Println("Processing name:", text)
		session.TempName = text
		session.State = StateKirimQty
		b.reply(msg.Chat.ID, "📏 Tovar miqdorini kiriting (masalan: 10):")

	case StateKirimQty:
		log.Println("Processing qty:", text) // ⬅️ Bu "name" emas "qty" bo'lishi kerak
		qty, err := strconv.ParseFloat(text, 64)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Noto'g'ri son formati. Iltimos, faqat raqam kiriting.")
			return
		}
		session.TempQty = qty
		session.State = StateKirimBuy
		b.reply(msg.Chat.ID, "💰 Kirim (olingan) narxini kiriting:")

	case StateKirimBuy:
		log.Println("Processing buy price:", text) // ⬅️ Tuzatildi
		price, err := strconv.ParseFloat(text, 64)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Noto'g'ri narx formati. Masalan: 25000")
			return
		}
		session.TempBuyPrice = price
		session.State = StateKirimSell
		b.reply(msg.Chat.ID, "🏷️ Sotuv narxini kiriting:")

	case StateKirimSell:
		log.Println("Processing sell price:", text) // ⬅️ Tuzatildi
		sell, err := strconv.ParseFloat(text, 64)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Noto'g'ri narx formati. Masalan: 40000")
			return
		}

		// ✅ Ma'lumotni bazaga saqlaymiz
		err = b.repo.AddProduct(context.Background(),
			session.TempName,
			session.TempQty,
			session.TempBuyPrice,
			sell,
		)
		if err != nil {
			b.reply(msg.Chat.ID, fmt.Sprintf("❌ Saqlashda xatolik: %v", err))
			delete(b.sessions, userID)
			return
		}

		b.reply(msg.Chat.ID,
			fmt.Sprintf("✅ %s dan %.2f dona qo'shildi!\n💰 Kirim: %.0f so'm\n🏷️ Sotuv: %.0f so'm",
				session.TempName, session.TempQty, session.TempBuyPrice, sell),
		)

		// 🔚 Sessionni tozalaymiz
		delete(b.sessions, userID)

	default:
		b.reply(msg.Chat.ID, "⚠️ Avval /kirim buyrug'ini yozing.")
	}
}

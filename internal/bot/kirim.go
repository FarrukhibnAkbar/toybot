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
	text := strings.TrimSpace(msg.Text)

	if !ok {
		b.sessions[userID] = &Session{State: StateKirimName}
		b.reply(msg.Chat.ID, "📦 Tovar nomini kiriting:")
		return
	}

	switch session.State {
	case StateKirimName:
		item, err := b.repo.GetItemByName(context.Background(), text)
		session.TempName = text
		if err == nil {
			// Mahsulot mavjud
			session.TempID = item.ID
			session.TempQty = item.Quantity
			session.TempBuyPrice = item.BuyPrice
			session.TempSellPrice = item.SellPrice
			session.State = StateKirimExistsChoice

			b.reply(msg.Chat.ID,
				fmt.Sprintf("⚠️ %s tovari mavjud:\nMiqdor: %.2f\nKirim narxi: %.0f\nSotuv narxi: %.0f\n\nYangilash uchun /edit kiriting.",
					item.Name, item.Quantity, item.BuyPrice, item.SellPrice))
			return
		}
		session.State = StateKirimQty
		b.reply(msg.Chat.ID, "📏 Miqdorni kiriting:")
		return

	case StateKirimExistsChoice:
		if text == "/new" {
			session.TempID = 0
			session.State = StateKirimQty
			b.reply(msg.Chat.ID, "📏 Miqdorni kiriting:")
			return
		}
		if text == "/edit" {
			session.State = StateKirimQty
			b.reply(msg.Chat.ID, "✏️ Yangi miqdorni kiriting:")
			log.Println("qty: ", session.TempQty)
			return
		}
		b.reply(msg.Chat.ID, "❌ /new yoki /edit kiriting.")
		return

	case StateKirimQty:
		qty, err := strconv.ParseFloat(text, 64)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Noto‘g‘ri miqdor.")
			return
		}
		log.Println("qty: ", qty)
		log.Println("tempqty: ", session.TempQty)
		session.TempQty = qty
		log.Println("updated tempqty: ", session.TempQty)
		session.State = StateKirimBuy
		b.reply(msg.Chat.ID, "💰 Kirim narxini kiriting:")
		return

	case StateKirimBuy:
		price, err := strconv.ParseFloat(text, 64)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Noto‘g‘ri narx.")
			return
		}
		session.TempBuyPrice = price
		session.State = StateKirimSell
		b.reply(msg.Chat.ID, "🏷️ Sotuv narxini kiriting:")
		return

	case StateKirimSell:
		sell, err := strconv.ParseFloat(text, 64)
		if err != nil {
			b.reply(msg.Chat.ID, "❌ Noto‘g‘ri narx.")
			return
		}
		session.TempSellPrice = sell

		user, _ := b.repo.GetUserDataByTgID(context.Background(), msg.From.ID)
		item, err := b.repo.GetItemByName(context.Background(), session.TempName)
		if session.TempID != 0 && err == nil {
			// 🔹 Eski qiymatlar
			oldQty := item.Quantity
			oldBuy := item.BuyPrice
			oldSell := item.SellPrice

			newQty := oldQty + session.TempQty
			newBuy := session.TempBuyPrice
			newSell := session.TempSellPrice

			_ = b.repo.CreateAuditLog(context.Background(), int(session.TempID),
				oldQty, oldBuy, oldSell,
				newQty, newBuy, newSell,
				int(user.ID),
			)
			log.Println("oldQty: ", oldQty)
			log.Println("tempQty + oldQty: ", session.TempQty+oldQty)
			log.Println("newQty: ", newQty)
			log.Println("tempQty: ", session.TempQty)
			_ = b.repo.UpdateItem(context.Background(), int(session.TempID), newQty, newBuy, newSell)

			b.reply(msg.Chat.ID, fmt.Sprintf("✏️ %s yangilandi!\nYangi miqdor: %.2f\nNarx: %.0f",
				session.TempName, newQty, newSell))
			delete(b.sessions, userID)
			return
		}

		// 🔸 Yangi mahsulot
		err = b.repo.AddProduct(context.Background(), session.TempName, session.TempQty, session.TempBuyPrice, session.TempSellPrice)
		if err != nil {
			b.reply(msg.Chat.ID, fmt.Sprintf("❌ Xatolik: %v", err))
			delete(b.sessions, userID)
			return
		}

		b.reply(msg.Chat.ID,
			fmt.Sprintf("✅ %s dan %.2f dona qo'shildi!\n💰 %.0f so‘m\n🏷️ %.0f so‘m",
				session.TempName, session.TempQty, session.TempBuyPrice, session.TempSellPrice),
		)
		delete(b.sessions, userID)
		return
	}
}

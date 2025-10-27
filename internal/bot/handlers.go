package bot

// func (b *Bot) handleKirim(msg *tgbotapi.Message, q *db.Queries) {
// 	args := strings.Fields(msg.Text)
// 	if len(args) < 5 {
// 		b.reply(msg.Chat.ID, "❗ Format: /kirim <nomi> <soni> <kelgan_narxi> <sotuv_narxi>")
// 		return
// 	}

// 	name := args[1]

// 	qty, err := strconv.ParseFloat(args[2], 64)
// 	if err != nil || qty <= 0 {
// 		b.reply(msg.Chat.ID, "❗ Noto'g'ri soni. Iltimos, musbat raqam kiriting.")
// 		return
// 	}

// 	costPrice, err := strconv.ParseFloat(args[3], 64)
// 	if err != nil || costPrice < 0 {
// 		b.reply(msg.Chat.ID, "❗ Noto'g'ri kelgan narxi. Iltimos, to'g'ri raqam kiriting.")
// 		return
// 	}

// 	sellPrice, err := strconv.ParseFloat(args[4], 64)
// 	if err != nil || sellPrice < 0 {
// 		b.reply(msg.Chat.ID, "❗ Noto'g'ri sotuv narxi. Iltimos, to'g'ri raqam kiriting.")
// 		return
// 	}

// 	ctx := context.Background()

// 	var product db.Product
// 	product, err = q.GetProductByName(ctx, name)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			product, err = q.CreateProduct(ctx, db.CreateProductParams{
// 				Name:      name,
// 				SellPrice: sellPrice,
// 			})
// 			if err != nil {
// 				b.reply(msg.Chat.ID, fmt.Sprintf("❌ Xatolik (product yaratishda): %v", err))
// 				return
// 			}
// 		} else {
// 			b.reply(msg.Chat.ID, fmt.Sprintf("❌ Xatolik (productni olishda): %v", err))
// 			return
// 		}
// 	}

// 	createdBy := int32(1) // default
// 	if msg.From != nil && msg.From.ID != 0 {
// 		createdBy = int32(msg.From.ID)
// 	}

// 	_, err = q.AddPurchase(ctx, db.AddPurchaseParams{
// 		ProductID: product.ID,
// 		Qty:       qty,
// 		CostPrice: costPrice,
// 		CreatedBy: createdBy,
// 	})
// 	if err != nil {
// 		b.reply(msg.Chat.ID, fmt.Sprintf("❌ Kirim yozishda xatolik: %v", err))
// 		return
// 	}

// 	if err := q.UpsertStock(ctx, db.UpsertStockParams{
// 		ProductID: int32(product.ID),
// 		Qty:       qty,
// 		CostPrice: costPrice,
// 	}); err != nil {
// 		b.reply(msg.Chat.ID, fmt.Sprintf("❌ Omborga yozishda xatolik: %v", err))
// 		return
// 	}

// 	b.reply(msg.Chat.ID, fmt.Sprintf("✅ %s dan %.2f dona qo‘shildi!", name, qty))
// }

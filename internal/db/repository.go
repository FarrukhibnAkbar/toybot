package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"toybot/internal/models"
)

// AddProduct — tovarni qo‘lda bazaga qo‘shish
func (q *Queries) AddProduct(ctx context.Context, name string, qty, buyPrice, sellPrice float64) error {
	_, err := q.db.Exec(
		ctx,
		`INSERT INTO products (name, quantity, purchase_price, sell_price)
         VALUES ($1, $2, $3, $4)`,
		name, qty, buyPrice, sellPrice,
	)
	if err != nil {
		log.Printf("❌ Xatolik (AddProduct): %v", err)
		return err
	}
	return nil
}

// 🔹 1) Mahsulotni nomi bilan olish
func (s *Queries) GetItemByName(ctx context.Context, name string) (*models.Item, error) {
	row := s.db.QueryRow(ctx, `
		SELECT id, name, quantity, purchase_price, sell_price
		FROM products
		WHERE name = $1
	`, name)

	var item models.Item
	if err := row.Scan(&item.ID, &item.Name, &item.Quantity, &item.BuyPrice, &item.SellPrice); err != nil {
		return nil, errors.New("mahsulot topilmadi")
	}
	return &item, nil
}

// 🔹 2) Mahsulot miqdorini yangilash
func (s *Queries) UpdateItemQuantity(ctx context.Context, id int, newQty float64) error {
	cmd, err := s.db.Exec(ctx, `
		UPDATE products SET quantity = $1 WHERE id = $2
	`, newQty, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("yangilash uchun mos mahsulot topilmadi")
	}
	return nil
}

// 🔹 3) Sotuv yozuvini yaratish
func (s *Queries) CreateSale(ctx context.Context, product_id int, qty, sellPrice float64, user_id int) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO sales (product_id, qty, sell_price, created_by)
		VALUES ($1, $2, $3, $4)
	`, product_id, qty, sellPrice, user_id)
	if err != nil {
		log.Printf("❌ Xatolik (CreateSale): %v", err)
		return fmt.Errorf("sotuv yozish xatosi: %v", err)
	}
	return nil
}

func (s *Queries) GetUserDataByTgID(ctx context.Context, tgID int64) (models.Users, error) {
	row := s.db.QueryRow(ctx, `
		SELECT id, full_name, role FROM users WHERE tg_id = $1
	`, tgID)

	var user models.Users
	if err := row.Scan(&user.ID, &user.Fullname, &user.Role); err != nil {
		log.Printf("❌ Xatolik (GetUserDataByTgID): %v", err)
		return models.Users{}, fmt.Errorf("foydalanuvchi topilmadi: %v", err)
	}
	return user, nil
}

func (s *Queries) CreateAuditLog(ctx context.Context, productID int, oldQty, oldBuy, oldSell, newQty, newBuy, newSell float64, userID int) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO audit_logs (
			product_id, old_quantity, old_buy_price, old_sell_price,
			new_quantity, new_buy_price, new_sell_price, changed_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, productID, oldQty, oldBuy, oldSell, newQty, newBuy, newSell, userID)
	return err
}

func (s *Queries) UpdateItem(ctx context.Context, id int, qty, buy, sell float64) error {
	_, err := s.db.Exec(ctx, `
		UPDATE products SET quantity = $1, purchase_price = $2, sell_price = $3 WHERE id = $4
	`, qty, buy, sell, id)
	return err
}
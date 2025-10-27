package db

import (
	"context"
	"log"
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

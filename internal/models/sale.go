package models

type Item struct {
	ID        int64
	Name      string
	Quantity  float64
	BuyPrice  float64
	SellPrice float64
}

type Sale struct {
	ID         int64
	Name       string
	Quantity   float64
	SellPrice  float64
	TotalPrice float64
	CreatedAt  string
}

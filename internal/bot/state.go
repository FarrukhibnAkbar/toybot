package bot

type UserState string

const (
	StateNone              UserState = ""
	StateKirimName         UserState = "kirim_name"
	StateKirimQty          UserState = "kirim_qty"
	StateKirimBuy          UserState = "kirim_buy"
	StateKirimSell         UserState = "kirim_sell"
	StateKirimExistsChoice UserState = "kirim_exists_choice"

	// yangi sotish holatlari
	StateSellName UserState = "sell_name"
	StateSellQty  UserState = "sell_qty"
)

type Session struct {
	State         UserState
	TempName      string
	TempQty       float64
	TempBuyPrice  float64
	TempSellPrice float64
	TempID        int64
}

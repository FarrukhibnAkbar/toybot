package bot

type UserState string

const (
	StateNone      UserState = ""
	StateKirimName UserState = "kirim_name"
	StateKirimQty  UserState = "kirim_qty"
	StateKirimBuy  UserState = "kirim_buy"
	StateKirimSell UserState = "kirim_sell"
)

type Session struct {
	State        UserState
	TempName     string
	TempQty      float64
	TempBuyPrice float64
}

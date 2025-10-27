package utils

import (
	"math/big"

	pgtype "github.com/jackc/pgx/v5/pgtype"
)

// FloatToNumeric converts float64 → pgx/v5 pgtype.Numeric
func FloatToNumeric(f float64) pgtype.Numeric {
	// float64 → string → big.Rat → Numeric
	rat := new(big.Rat).SetFloat64(f)
	num := pgtype.Numeric{}
	num.Int = rat.Num()
	num.Exp = int32(-len(rat.Denom().String()) + 1)
	// num.Status = pgtype.Present
	return num
}

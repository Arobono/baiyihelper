package cal

import "github.com/shopspring/decimal"

//加
func Add(a ...float64) float64 {
	var dec = decimal.NewFromFloat(0)
	for _, i := range a {
		dec = dec.Add(decimal.NewFromFloat(i))
	}
	res, _ := dec.Float64()
	return res
}

//除
func Divide(a, b float64) float64 {
	da := decimal.NewFromFloat(a)
	db := decimal.NewFromFloat(b)
	res, _ := da.Div(db).Float64()
	return res
}

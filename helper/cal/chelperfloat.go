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

//减
func Subtract(a, b float64) float64 {
	da := decimal.NewFromFloat(a)
	db := decimal.NewFromFloat(b)
	res, _ := da.Sub(db).Float64()
	return res
}

//乘
func Multiply(a ...float64) float64 {
	var dec = decimal.NewFromFloat(1)
	for _, i := range a {
		dec = dec.Mul(decimal.NewFromFloat(i))
	}
	res, _ := dec.Float64()
	return res

	// da := decimal.NewFromFloat(a)
	// db := decimal.NewFromFloat(b)
	// res, _ := da.Mul(db).Float64()
	// return res
}

//位数限制
func Digits(a float64, mod int) float64 {
	decimal.DivisionPrecision = mod
	res, _ := decimal.NewFromFloat(a).Float64()
	return res
}

//float64  a是否在b、c之间
func Between(da, db, dc float64) bool {
	if db == dc {
		return da == db
	}
	var maxv float64
	var minv float64
	if db > dc {
		maxv = db
		minv = dc
	} else {
		maxv = dc
		minv = db
	}

	return (maxv >= da && da >= minv)
}

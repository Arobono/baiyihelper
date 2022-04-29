package mode

import "time"

type PositionInput struct {
	RealAccountId string
	DateTime      time.Time
	Symbol        string
	SymbolName    string
	Direction     string
	EntrustPrice  float64
	EntrustHand   float64
}

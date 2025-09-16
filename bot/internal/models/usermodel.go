package models

type User struct {
	ID            int
	ChatID        int64
	Username      string
	Balance       int
	AutoBuy       bool
	MinCostLimit  int
	MaxCostLimit  int
	CyclesCount   int
	SupplyLimit   int
	ChannelEnabled bool
}
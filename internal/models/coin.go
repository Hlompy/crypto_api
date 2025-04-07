package models

type UserCoin struct {
	ID         int     `json:"id,omitempty"`
	Name       string  `json:"name" binding:"required"`
	Symbol     string  `json:"symbol" binding:"required"`
	UsdPrice   float64 `json:"usd_price" binding:"required"`
	IsUserCoin bool    `json:"is_user_coin,omitempty"`
}

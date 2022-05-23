package dto

import (
	"time"
)

// PlaceWagerRequest ...
type PlaceWagerRequest struct {
	TotalWagerValue   uint32  `json:"total_wager_value"`
	Odds              uint32  `json:"odds"`
	SellingPercentage float32 `json:"selling_percentage"`
	SellingPrice      float32 `json:"selling_price"`
}

// Wager ...
type Wager struct {
	ID                  uint32     `json:"id"`
	TotalWagerValue     uint32     `json:"total_wager_value"`
	Odds                uint32     `json:"odds"`
	SellingPercentage   float32    `json:"selling_percentage"`
	SellingPrice        float32    `json:"selling_price"`
	CurrentSellingPrice float32    `json:"current_selling_price"`
	PercentageSold      float32    `json:"percentage_sold"`
	AmountSold          float32    `json:"amount_sold"`
	PlacedAt            *time.Time `json:"placed_at"`
}

// BuyWagerRequest ...
type BuyWagerRequest struct {
	BuyingPrice float32 `json:"buying_price"`
}

// ListWagerRequest ...
type ListWagerRequest struct {
	Page  uint32
	Limit uint32
}

// WagerPurchase ...
type WagerPurchase struct {
	ID          uint32     `json:"id"`
	WagerID     uint32     `json:"wager_id"`
	BuyingPrice float32    `json:"buying_Price"`
	BoughtAt    *time.Time `json:"bought_at"`
}

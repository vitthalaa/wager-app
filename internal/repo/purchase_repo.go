package repo

import (
	"context"
	"database/sql"
)

// Purchase ...
type Purchase struct {
	ID          uint32
	WagerID     uint32
	BuyingPrice float32
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

// IPurchaseRepo is repository interface for purchase db operations
type IPurchaseRepo interface {
	CreatePurchase(ctx context.Context, wager *Purchase) (*Purchase, error)
}

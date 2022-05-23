package repo

import (
	"context"
	"database/sql"
)

const (
	insertWagerStmt = `insert into wager(total_wager_value, odds, selling_percentage, selling_price, current_selling_price)
						values ($1, $2, $3, $4, $5)
						returning *`
)

// Wager ...
type Wager struct {
	ID                  uint32
	TotalWagerValue     uint32
	Odds                uint32
	SellingPercentage   float32
	SellingPrice        float32
	CurrentSellingPrice float32
	PercentageSold      sql.NullFloat64
	AmountSold          sql.NullFloat64
	CreatedAt           sql.NullTime
	UpdatedAt           sql.NullTime
}

// IWagerRepo is repository interface for wager db operations
type IWagerRepo interface {
	CreateWager(ctx context.Context, wager *Wager) (*Wager, error)
}

// NewWagerRepo ...
func NewWagerRepo(db *sql.DB) *WagerRepo {
	return &WagerRepo{
		db: db,
	}
}

// WagerRepo is repository implementation for wager db operations
type WagerRepo struct {
	db *sql.DB
}

// CreateWager creates new wager record in db
func (wr *WagerRepo) CreateWager(ctx context.Context, wager *Wager) (*Wager, error) {
	row := wr.db.QueryRowContext(ctx,
		insertWagerStmt,
		wager.TotalWagerValue, wager.Odds, wager.SellingPercentage, wager.SellingPrice, wager.CurrentSellingPrice)

	err := row.Scan(
		&wager.ID,
		&wager.TotalWagerValue,
		&wager.Odds,
		&wager.SellingPercentage,
		&wager.SellingPrice,
		&wager.CurrentSellingPrice,
		&wager.PercentageSold,
		&wager.AmountSold,
		&wager.CreatedAt,
		&wager.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return wager, nil
}

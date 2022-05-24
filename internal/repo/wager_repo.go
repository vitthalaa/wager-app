package repo

import (
	"context"
	"database/sql"
)

const (
	insertWagerStmt = `insert into wager(total_wager_value, odds, selling_percentage, selling_price, current_selling_price)
						values ($1, $2, $3, $4, $5)
						returning *`
	listWagerStmt    = "select * from wager order by id desc limit $1 offset $2"
	getWagerByIDStmt = "select * from wager where id=$1"
	updateWagerStmt  = `update wager set current_selling_price=$1, percentage_sold=$2, amount_sold=$3, updated_at=now() 
						where id = $4;`
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
	AmountSold          sql.NullInt32
	CreatedAt           sql.NullTime
	UpdatedAt           sql.NullTime
}

// IWagerRepo is repository interface for wager db operations
type IWagerRepo interface {
	CreateWager(ctx context.Context, wager *Wager) (*Wager, error)
	ListWager(ctx context.Context, offset, limit uint32) ([]Wager, error)
	GetWagerByID(ctx context.Context, wagerID uint32) (*Wager, error)
	UpdateWager(ctx context.Context, wager *Wager) error
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
	stmt, err := wr.db.Prepare(insertWagerStmt)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(ctx,
		wager.TotalWagerValue, wager.Odds, wager.SellingPercentage, wager.SellingPrice, wager.CurrentSellingPrice)

	err = row.Scan(
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

// ListWager returns list of wagers from offset to limit
func (wr *WagerRepo) ListWager(ctx context.Context, offset, limit uint32) ([]Wager, error) {
	stmt, err := wr.db.Prepare(listWagerStmt)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]Wager, 0, limit)
	for rows.Next() {
		var wager Wager
		err = rows.Scan(
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

		res = append(res, wager)
	}

	return res, nil
}

// GetWagerByID returns wager record by ids
func (wr *WagerRepo) GetWagerByID(ctx context.Context, wagerID uint32) (*Wager, error) {
	stmt, err := wr.db.Prepare(getWagerByIDStmt)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, wagerID)

	var wager Wager
	err = row.Scan(
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

	return &wager, nil
}

// UpdateWager updates wager record for current selling price, amount sold and percentage sold
func (wr *WagerRepo) UpdateWager(ctx context.Context, wager *Wager) error {
	stmt, err := wr.db.Prepare(updateWagerStmt)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, wager.CurrentSellingPrice, wager.PercentageSold, wager.AmountSold, wager.ID)
	if err != nil {
		return err
	}

	return nil
}

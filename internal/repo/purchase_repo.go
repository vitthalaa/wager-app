package repo

import (
	"context"
	"database/sql"
)

const (
	insertPurchaseStmt = `insert into purchases(wager_id, buying_price) values ($1, $2)
						returning *`
	deletePurchaseStmt = "delete from purchases where id = $1"
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
	CreatePurchase(ctx context.Context, purchase *Purchase) (*Purchase, error)
	DeletePurchase(ctx context.Context, id uint32) error
}

// NewPurchaseRepo ...
func NewPurchaseRepo(db *sql.DB) *PurchaseRepo {
	return &PurchaseRepo{
		db: db,
	}
}

// PurchaseRepo is repository implementation for wager db operations
type PurchaseRepo struct {
	db *sql.DB
}

// CreatePurchase creates new purchase record in db
func (pr *PurchaseRepo) CreatePurchase(ctx context.Context, purchase *Purchase) (*Purchase, error) {
	stmt, err := pr.db.Prepare(insertPurchaseStmt)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(ctx,
		purchase.WagerID, purchase.BuyingPrice)

	err = row.Scan(
		&purchase.ID,
		&purchase.WagerID,
		&purchase.BuyingPrice,
		&purchase.CreatedAt,
		&purchase.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return purchase, nil
}

// DeletePurchase deletes purchase record from db by id
func (pr *PurchaseRepo) DeletePurchase(ctx context.Context, id uint32) error {
	stmt, err := pr.db.Prepare(deletePurchaseStmt)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

package repo

import (
	"context"
	"database/sql"
)

// IWagerRepo is repository interface for wager db operations
type IWagerRepo interface {
	CreateWager(ctx context.Context) error
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
func (wr *WagerRepo) CreateWager(ctx context.Context) error {
	panic("not implemented")
}

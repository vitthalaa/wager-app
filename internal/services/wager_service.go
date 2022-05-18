package services

import (
	"context"

	"github.com/vitthalaa/wager-app/internal/repo"
)

// IWagerService ...
type IWagerService interface {
	PlaceWager(ctx context.Context) error
	BuyWager(ctx context.Context) error
	ListWager(ctx context.Context) error
}

// NewWagerService ...
func NewWagerService(wagerRepo repo.IWagerRepo) *WagerService {
	return &WagerService{
		wagerRepo: wagerRepo,
	}
}

// WagerService ...
type WagerService struct {
	wagerRepo repo.IWagerRepo
}

// PlaceWager ...
func (s *WagerService) PlaceWager(ctx context.Context) error {
	panic("not implemented")
}

// BuyWager ...
func (s *WagerService) BuyWager(ctx context.Context) error {
	panic("not implemented")
}

// ListWager ...
func (s *WagerService) ListWager(ctx context.Context) error {
	panic("not implemented")
}

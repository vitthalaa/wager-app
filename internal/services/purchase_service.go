package services

import (
	"context"

	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/repo"
)

type IPurchaseService interface {
	BuyWager(ctx context.Context, req *dto.BuyWagerRequest) (*dto.WagerPurchase, error)
}

// NewPurchaseService ...
func NewPurchaseService(purchaseRepo repo.IPurchaseRepo) *PurchaseService {
	return &PurchaseService{
		purchaseRepo: purchaseRepo,
	}
}

// PurchaseService ...
type PurchaseService struct {
	purchaseRepo repo.IPurchaseRepo
}

func (s *PurchaseService) BuyWager(ctx context.Context, req *dto.BuyWagerRequest) (*dto.WagerPurchase, error) {
	panic("not implemented")
}

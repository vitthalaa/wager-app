package services

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/vitthalaa/wager-app/app_errors"
	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/repo"
)

// IPurchaseService ...
type IPurchaseService interface {
	PurchaseWager(ctx context.Context, req *dto.BuyWagerRequest) (*dto.WagerPurchase, error)
}

// NewPurchaseService ...
func NewPurchaseService(purchaseRepo repo.IPurchaseRepo, wagerRepo repo.IWagerRepo) *PurchaseService {
	return &PurchaseService{
		purchaseRepo: purchaseRepo,
		wagerRepo:    wagerRepo,
	}
}

// PurchaseService ...
type PurchaseService struct {
	purchaseRepo repo.IPurchaseRepo
	wagerRepo    repo.IWagerRepo
}

// PurchaseWager ...
func (s *PurchaseService) PurchaseWager(ctx context.Context, req *dto.BuyWagerRequest) (*dto.WagerPurchase, error) {
	if req == nil || req.WagerID == 0 {
		return nil, &app_errors.ErrorResponse{Status: http.StatusBadRequest, Code: app_errors.ErrInvalidWagerID}
	}

	if req.BuyingPrice < 1 {
		return nil, &app_errors.ErrorResponse{Status: http.StatusBadRequest, Code: app_errors.ErrInvalidBuyingPrice}
	}

	// TODO: Should do all wager operations with lock and transaction
	wager, err := s.wagerRepo.GetWagerByID(ctx, req.WagerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &app_errors.ErrorResponse{Status: http.StatusNotFound, Code: app_errors.ErrNotFound}
		}

		return nil, err
	}

	if req.BuyingPrice > wager.CurrentSellingPrice {
		return nil, &app_errors.ErrorResponse{Status: http.StatusBadRequest, Code: app_errors.ErrInvalidBuyingPrice}
	}

	// TODO: Clarify whether need to return error if amount sold is reaches to total wager value
	// Or percent sold reaches to selling percent
	if wager.TotalWagerValue <= uint32(wager.AmountSold.Int32) {
		return nil, &app_errors.ErrorResponse{Status: http.StatusNotAcceptable, Code: app_errors.ErrWagerSoldOut}
	}

	// insert purchase
	purchaseReq := &repo.Purchase{
		WagerID:     wager.ID,
		BuyingPrice: req.BuyingPrice,
	}

	purchase, err := s.purchaseRepo.CreatePurchase(ctx, purchaseReq)
	if err != nil {
		return nil, err
	}

	// increase amount sold
	wager.AmountSold = sql.NullInt32{
		Int32: wager.AmountSold.Int32 + 1,
		Valid: true,
	}

	// buying price will be assigned to current selling price
	wager.CurrentSellingPrice = req.BuyingPrice
	wager.PercentageSold = sql.NullFloat64{
		Float64: float64(wager.AmountSold.Int32*100) / float64(wager.TotalWagerValue),
		Valid:   true,
	}

	err = s.wagerRepo.UpdateWager(ctx, wager)
	if err != nil {
		ctxBg, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		go s.revertPurchase(ctxBg, cancel, purchase.ID)
		return nil, err
	}

	// return purchase
	purchaseDTO := &dto.WagerPurchase{
		ID:          purchase.ID,
		WagerID:     purchase.WagerID,
		BuyingPrice: purchase.BuyingPrice,
	}

	if purchase.CreatedAt.Valid {
		boughtAt := purchase.CreatedAt.Time
		purchaseDTO.BoughtAt = &boughtAt
	}

	return purchaseDTO, nil
}

func (s *PurchaseService) revertPurchase(ctx context.Context, cancel context.CancelFunc, id uint32) {
	err := s.purchaseRepo.DeletePurchase(ctx, id)
	if err != nil {
		log.Printf("delete purchase error %s", err)
		return
	}

	log.Printf("purchase record %d deleted", id)
	cancel()
}

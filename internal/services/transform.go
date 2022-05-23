package services

import (
	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/repo"
)

func toWagerEntity(req dto.PlaceWagerRequest) *repo.Wager {
	return &repo.Wager{
		TotalWagerValue:     req.TotalWagerValue,
		Odds:                req.Odds,
		SellingPercentage:   req.SellingPercentage,
		SellingPrice:        req.SellingPrice,
		CurrentSellingPrice: req.SellingPrice,
	}
}

func toWagerDTO(w *repo.Wager) *dto.Wager {
	if w == nil {
		return nil
	}

	wDto := &dto.Wager{
		ID:                  w.ID,
		TotalWagerValue:     w.TotalWagerValue,
		Odds:                w.Odds,
		SellingPercentage:   w.SellingPercentage,
		SellingPrice:        w.SellingPrice,
		CurrentSellingPrice: w.CurrentSellingPrice,
		PercentageSold:      float32(w.PercentageSold.Float64),
		AmountSold:          float32(w.AmountSold.Float64),
	}

	if w.CreatedAt.Valid {
		t := w.CreatedAt.Time
		wDto.PlacedAt = &t
	}

	return wDto
}
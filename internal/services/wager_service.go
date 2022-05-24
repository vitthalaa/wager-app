package services

import (
	"context"

	"github.com/vitthalaa/wager-app/app_errors"
	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/repo"
)

// IWagerService ...
type IWagerService interface {
	PlaceWager(ctx context.Context, req *dto.PlaceWagerRequest) (*dto.Wager, error)
	ListWager(ctx context.Context, req *dto.ListWagerRequest) ([]dto.Wager, error)
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
func (s *WagerService) PlaceWager(ctx context.Context, req *dto.PlaceWagerRequest) (*dto.Wager, error) {
	errRes := validatePlaceWagerRequest(req)
	if errRes != nil {
		return nil, errRes
	}

	wager, err := s.wagerRepo.CreateWager(ctx, toWagerEntity(*req))
	if err != nil {
		return nil, err
	}

	wagerDto := toWagerDTO(*wager)
	return &wagerDto, nil
}

// ListWager ...
func (s *WagerService) ListWager(ctx context.Context, req *dto.ListWagerRequest) ([]dto.Wager, error) {
	offset := uint32(0)
	limit := uint32(10)
	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.Page > 0 {
		offset = (req.Page - 1) * limit
	}

	res, err := s.wagerRepo.ListWager(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	dtoList := make([]dto.Wager, 0, len(res))
	for _, rs := range res {
		dtoList = append(dtoList, toWagerDTO(rs))
	}

	return dtoList, nil
}

func validatePlaceWagerRequest(req *dto.PlaceWagerRequest) *app_errors.ErrorResponse {
	err := &app_errors.ErrorResponse{
		Status: 400,
	}
	switch true {
	case req.TotalWagerValue < 1:
		err.Code = app_errors.ErrInvalidTotalWagerValue
		return err

	case req.Odds < 1:
		err.Code = app_errors.ErrInvalidOdds
		return err

	case req.SellingPercentage < 1 || req.SellingPercentage > 100:
		err.Code = app_errors.ErrInvalidSellingPercentage
		return err

	case req.SellingPrice <= float32(req.TotalWagerValue)*(req.SellingPercentage/100):
		err.Code = app_errors.ErrInvalidSellingPrice
		return err
	}

	return nil
}

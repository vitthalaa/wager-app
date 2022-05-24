package services

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vitthalaa/wager-app/app_errors"
	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/repo"
)

func TestPurchaseService_PurchaseWager(t *testing.T) {
	now := time.Now()
	for _, tc := range []struct {
		name  string
		input *dto.BuyWagerRequest

		wagerRepoResp  *repo.Wager
		wagerRepoError error

		purchaseRepoResp  *repo.Purchase
		purchaseRepoError error

		updateWagerRepoReq   *repo.Wager
		updateWagerRepoError error

		expectedRes   *dto.WagerPurchase
		expectedError error
	}{
		{
			name: "happy path",
			input: &dto.BuyWagerRequest{
				WagerID:     111,
				BuyingPrice: 25.5,
			},
			wagerRepoResp: &repo.Wager{
				ID:                  111,
				TotalWagerValue:     100,
				Odds:                2,
				SellingPercentage:   20,
				SellingPrice:        26,
				CurrentSellingPrice: 26,
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			wagerRepoError: nil,
			purchaseRepoResp: &repo.Purchase{
				ID:          1,
				WagerID:     111,
				BuyingPrice: 25.5,
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			purchaseRepoError: nil,
			updateWagerRepoReq: &repo.Wager{
				ID:                  111,
				TotalWagerValue:     100,
				Odds:                2,
				SellingPercentage:   20,
				SellingPrice:        26,
				CurrentSellingPrice: 25.5,
				PercentageSold: sql.NullFloat64{
					Float64: 1,
					Valid:   true,
				},
				AmountSold: sql.NullInt32{
					Int32: 1,
					Valid: true,
				},
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			updateWagerRepoError: nil,
			expectedRes: &dto.WagerPurchase{
				ID:          1,
				WagerID:     111,
				BuyingPrice: 25.5,
				BoughtAt:    &now,
			},
			expectedError: nil,
		},
		{
			name:                 "invalid request",
			input:                nil,
			wagerRepoResp:        nil,
			wagerRepoError:       nil,
			purchaseRepoResp:     nil,
			purchaseRepoError:    nil,
			updateWagerRepoReq:   nil,
			updateWagerRepoError: nil,
			expectedRes:          nil,
			expectedError:        &app_errors.ErrorResponse{Status: http.StatusBadRequest, Code: app_errors.ErrInvalidWagerID},
		},
		{
			name:                 "invalid wager id",
			input:                &dto.BuyWagerRequest{},
			wagerRepoResp:        nil,
			wagerRepoError:       nil,
			purchaseRepoResp:     nil,
			purchaseRepoError:    nil,
			updateWagerRepoReq:   nil,
			updateWagerRepoError: nil,
			expectedRes:          nil,
			expectedError:        &app_errors.ErrorResponse{Status: http.StatusBadRequest, Code: app_errors.ErrInvalidWagerID},
		},
		{
			name: "invalid buying price",
			input: &dto.BuyWagerRequest{
				WagerID:     111,
				BuyingPrice: 0,
			},
			wagerRepoResp:        nil,
			wagerRepoError:       nil,
			purchaseRepoResp:     nil,
			purchaseRepoError:    nil,
			updateWagerRepoReq:   nil,
			updateWagerRepoError: nil,
			expectedRes:          nil,
			expectedError:        &app_errors.ErrorResponse{Status: http.StatusBadRequest, Code: app_errors.ErrInvalidBuyingPrice},
		},
		{
			name: "get wager repo not found error",
			input: &dto.BuyWagerRequest{
				WagerID:     111,
				BuyingPrice: 22,
			},
			wagerRepoResp:        nil,
			wagerRepoError:       sql.ErrNoRows,
			purchaseRepoResp:     nil,
			purchaseRepoError:    nil,
			updateWagerRepoReq:   nil,
			updateWagerRepoError: nil,
			expectedRes:          nil,
			expectedError:        &app_errors.ErrorResponse{Status: http.StatusNotFound, Code: app_errors.ErrNotFound},
		},
		{
			name: "get wager repo unknown error",
			input: &dto.BuyWagerRequest{
				WagerID:     111,
				BuyingPrice: 22,
			},
			wagerRepoResp:        nil,
			wagerRepoError:       errors.New("some repo error"),
			purchaseRepoResp:     nil,
			purchaseRepoError:    nil,
			updateWagerRepoReq:   nil,
			updateWagerRepoError: nil,
			expectedRes:          nil,
			expectedError:        errors.New("some repo error"),
		},
		{
			name: "request buying price is greater than current selling error",
			input: &dto.BuyWagerRequest{
				WagerID:     111,
				BuyingPrice: 27,
			},
			wagerRepoResp: &repo.Wager{
				ID:                  111,
				TotalWagerValue:     100,
				Odds:                2,
				SellingPercentage:   20,
				SellingPrice:        26,
				CurrentSellingPrice: 26,
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			wagerRepoError:       nil,
			purchaseRepoResp:     nil,
			purchaseRepoError:    nil,
			updateWagerRepoReq:   nil,
			updateWagerRepoError: nil,
			expectedRes:          nil,
			expectedError:        &app_errors.ErrorResponse{Status: http.StatusBadRequest, Code: app_errors.ErrInvalidBuyingPrice},
		},
		{
			name: "purchase repo error",
			input: &dto.BuyWagerRequest{
				WagerID:     111,
				BuyingPrice: 25.5,
			},
			wagerRepoResp: &repo.Wager{
				ID:                  111,
				TotalWagerValue:     100,
				Odds:                2,
				SellingPercentage:   20,
				SellingPrice:        26,
				CurrentSellingPrice: 26,
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			wagerRepoError:       nil,
			purchaseRepoResp:     nil,
			purchaseRepoError:    errors.New("some purchase repo error"),
			updateWagerRepoReq:   nil,
			updateWagerRepoError: nil,
			expectedRes:          nil,
			expectedError:        errors.New("some purchase repo error"),
		},
		{
			name: "update wager error",
			input: &dto.BuyWagerRequest{
				WagerID:     111,
				BuyingPrice: 25.5,
			},
			wagerRepoResp: &repo.Wager{
				ID:                  111,
				TotalWagerValue:     100,
				Odds:                2,
				SellingPercentage:   20,
				SellingPrice:        26,
				CurrentSellingPrice: 26,
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			wagerRepoError: nil,
			purchaseRepoResp: &repo.Purchase{
				ID:          1,
				WagerID:     111,
				BuyingPrice: 25.5,
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			purchaseRepoError: nil,
			updateWagerRepoReq: &repo.Wager{
				ID:                  111,
				TotalWagerValue:     100,
				Odds:                2,
				SellingPercentage:   20,
				SellingPrice:        26,
				CurrentSellingPrice: 25.5,
				PercentageSold: sql.NullFloat64{
					Float64: 1,
					Valid:   true,
				},
				AmountSold: sql.NullInt32{
					Int32: 1,
					Valid: true,
				},
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			updateWagerRepoError: errors.New("some update wager repo error"),
			expectedRes:          nil,
			expectedError:        errors.New("some update wager repo error"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockWagerRepo := new(MockWagerRepo)
			if tc.input != nil {
				mockWagerRepo.On("GetWagerByID", ctx, tc.input.WagerID).
					Return(tc.wagerRepoResp, tc.wagerRepoError)
			}

			mockPurchaseRepo := new(MockPurchaseRepo)
			mockPurchaseRepo.On("CreatePurchase", ctx, mock.Anything).
				Return(tc.purchaseRepoResp, tc.purchaseRepoError)

			if tc.input != nil && tc.input.WagerID != 0 {
				mockPurchaseRepo.On("DeletePurchase", mock.Anything, mock.Anything).
					Return(nil)
			}

			mockWagerRepo.On("UpdateWager", ctx, tc.updateWagerRepoReq).
				Return(tc.updateWagerRepoError)

			service := NewPurchaseService(mockPurchaseRepo, mockWagerRepo)

			wagerPurchase, err := service.PurchaseWager(ctx, tc.input)

			assert.Equal(t, tc.expectedError, err)
			assert.EqualValues(t, tc.expectedRes, wagerPurchase)

		})
	}
}

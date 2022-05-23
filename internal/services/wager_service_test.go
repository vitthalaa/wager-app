package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vitthalaa/wager-app/app_errors"
	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/repo"
)

func TestWagerService_PlaceWager(t *testing.T) {
	now := time.Now()
	for _, tc := range []struct {
		name  string
		input *dto.PlaceWagerRequest

		repoResp  *repo.Wager
		repoError error

		expectedRes   *dto.Wager
		expectedError error
	}{
		{
			name: "happy path",
			input: &dto.PlaceWagerRequest{
				TotalWagerValue:   1000,
				Odds:              2,
				SellingPercentage: 20,
				SellingPrice:      201,
			},
			repoResp: &repo.Wager{
				ID:                  111,
				TotalWagerValue:     1000,
				Odds:                2,
				SellingPercentage:   20,
				SellingPrice:        201,
				CurrentSellingPrice: 201,
				CreatedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			repoError: nil,
			expectedRes: &dto.Wager{
				ID:                  111,
				TotalWagerValue:     1000,
				Odds:                2,
				SellingPercentage:   20,
				SellingPrice:        201,
				CurrentSellingPrice: 201,
				PercentageSold:      0,
				AmountSold:          0,
				PlacedAt:            &now,
			},
			expectedError: nil,
		},
		{
			name: "validation error",
			input: &dto.PlaceWagerRequest{
				TotalWagerValue:   0,
				Odds:              2,
				SellingPercentage: 20,
				SellingPrice:      201,
			},
			repoResp:    nil,
			repoError:   nil,
			expectedRes: nil,
			expectedError: &app_errors.ErrorResponse{
				Status: 400,
				Code:   app_errors.ErrInvalidTotalWagerValue,
			},
		},
		{
			name: "wager repo error",
			input: &dto.PlaceWagerRequest{
				TotalWagerValue:   1000,
				Odds:              2,
				SellingPercentage: 20,
				SellingPrice:      201,
			},
			repoResp:      nil,
			repoError:     errors.New("some repo error"),
			expectedRes:   nil,
			expectedError: errors.New("some repo error"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockWagerRepo)
			mockRepo.On("CreateWager", ctx, mock.Anything).
				Return(tc.repoResp, tc.repoError)

			service := NewWagerService(mockRepo)

			wager, err := service.PlaceWager(ctx, tc.input)

			assert.Equal(t, tc.expectedRes, wager)
			assert.Equal(t, err, tc.expectedError)
		})
	}
}

func Test_validatePlaceWagerRequest(t *testing.T) {
	for _, tc := range []struct {
		name          string
		req           *dto.PlaceWagerRequest
		expectedError *app_errors.ErrorResponse
	}{
		{
			name: "happy path",
			req: &dto.PlaceWagerRequest{
				TotalWagerValue:   1000,
				Odds:              2,
				SellingPercentage: 20,
				SellingPrice:      201,
			},
			expectedError: nil,
		},
		{
			name: "invalid total wager value",
			req: &dto.PlaceWagerRequest{
				TotalWagerValue:   0,
				Odds:              2,
				SellingPercentage: 20,
				SellingPrice:      201,
			},
			expectedError: &app_errors.ErrorResponse{
				Status: 400,
				Code:   app_errors.ErrInvalidTotalWagerValue,
			},
		},
		{
			name: "invalid odds",
			req: &dto.PlaceWagerRequest{
				TotalWagerValue:   1000,
				Odds:              0,
				SellingPercentage: 20,
				SellingPrice:      201,
			},
			expectedError: &app_errors.ErrorResponse{
				Status: 400,
				Code:   app_errors.ErrInvalidOdds,
			},
		},
		{
			name: "invalid selling percentage",
			req: &dto.PlaceWagerRequest{
				TotalWagerValue:   1000,
				Odds:              2,
				SellingPercentage: 101,
				SellingPrice:      201,
			},
			expectedError: &app_errors.ErrorResponse{
				Status: 400,
				Code:   app_errors.ErrInvalidSellingPercentage,
			},
		},
		{
			name: "invalid selling price",
			req: &dto.PlaceWagerRequest{
				TotalWagerValue:   1000,
				Odds:              2,
				SellingPercentage: 20,
				SellingPrice:      100,
			},
			expectedError: &app_errors.ErrorResponse{
				Status: 400,
				Code:   app_errors.ErrInvalidSellingPrice,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePlaceWagerRequest(tc.req)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

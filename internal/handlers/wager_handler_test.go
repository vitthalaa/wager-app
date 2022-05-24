package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/vitthalaa/wager-app/dto"
)

func TestWagersHandler_Handle_PlaceWager_HappyPath(t *testing.T) {
	placeWagerReq := dto.PlaceWagerRequest{
		TotalWagerValue:   1000,
		Odds:              2,
		SellingPercentage: 20,
		SellingPrice:      201,
	}

	now := time.Now()
	placeWagerRes := &dto.Wager{
		ID:                  1000,
		TotalWagerValue:     1000,
		Odds:                2,
		SellingPercentage:   20,
		SellingPrice:        201,
		CurrentSellingPrice: 201,
		PercentageSold:      0,
		AmountSold:          0,
		PlacedAt:            &now,
	}

	body, err := json.Marshal(placeWagerReq)
	require.Nil(t, err)

	request, err := http.NewRequest("POST", "/wagers", bytes.NewReader(body))
	require.Nil(t, err)

	mockWagerService := new(MockWagerService)
	mockWagerService.On("PlaceWager", mock.Anything, &placeWagerReq).
		Return(placeWagerRes, nil)

	resRecorder := httptest.NewRecorder()
	handler := NewWagersHandler(mockWagerService)
	handler.Handle(resRecorder, request)

	expected, err := json.Marshal(placeWagerRes)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resRecorder.Code)
	assert.Equal(t, string(expected), resRecorder.Body.String())
}

func TestWagersHandler_Handle_ListWager_HappyPath(t *testing.T) {
	now := time.Now()
	wagerListResp := []dto.Wager{
		{
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
		{
			ID:                  222,
			TotalWagerValue:     100,
			Odds:                2,
			SellingPercentage:   20.5,
			SellingPrice:        20.6,
			CurrentSellingPrice: 40.6,
			PercentageSold:      50.1,
			AmountSold:          50.1,
			PlacedAt:            &now,
		},
	}

	request, err := http.NewRequest("GET", "http://domain.co/wagers?page=1&limit=20", nil)
	require.Nil(t, err)

	expectedRequest := &dto.ListWagerRequest{
		Page:  uint32(1),
		Limit: uint32(20),
	}

	mockWagerService := new(MockWagerService)
	mockWagerService.On("ListWager", mock.Anything, expectedRequest).
		Return(wagerListResp, nil)

	resRecorder := httptest.NewRecorder()
	handler := NewWagersHandler(mockWagerService)
	handler.Handle(resRecorder, request)

	expected, err := json.Marshal(wagerListResp)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resRecorder.Code)
	assert.Equal(t, string(expected), resRecorder.Body.String())
}

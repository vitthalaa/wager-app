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

func TestWagersHandler_Handle_HappyPath(t *testing.T) {
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

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

func TestPurchaseHandler_Handle(t *testing.T) {
	buyWagerReq := dto.BuyWagerRequest{
		BuyingPrice: 25,
	}

	now := time.Now()
	purchaseRes := dto.WagerPurchase{
		ID:          1,
		WagerID:     111,
		BuyingPrice: 25,
		BoughtAt:    &now,
	}

	body, err := json.Marshal(buyWagerReq)
	require.Nil(t, err)

	request, err := http.NewRequest("POST", "http://domain.co/buy/111", bytes.NewReader(body))
	require.Nil(t, err)

	expectedReq := &dto.BuyWagerRequest{
		WagerID:     111,
		BuyingPrice: 25,
	}

	mockPurchaseService := new(MockPurchaseService)
	mockPurchaseService.On("PurchaseWager", mock.Anything, expectedReq).
		Return(&purchaseRes, nil)

	resRecorder := httptest.NewRecorder()
	handler := NewPurchasesHandler(mockPurchaseService)
	handler.Handle(resRecorder, request)

	expected, err := json.Marshal(purchaseRes)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resRecorder.Code)
	assert.Equal(t, string(expected), resRecorder.Body.String())
}

//go:build integration
// +build integration

package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	env "github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/config"
	"github.com/vitthalaa/wager-app/internal/db"
	"github.com/vitthalaa/wager-app/internal/handlers"
	"github.com/vitthalaa/wager-app/internal/repo"
	"github.com/vitthalaa/wager-app/internal/services"
)

func Test_WagerHandler(t *testing.T) {
	// 1. Create Wager
	placeWagerReq := dto.PlaceWagerRequest{
		TotalWagerValue:   100,
		Odds:              2,
		SellingPercentage: 20,
		SellingPrice:      21,
	}

	var loadEnv = env.Overload
	err := loadEnv("../.env")

	require.Nil(t, err)

	conf := config.GetAppConfig()

	conn, err := db.OpenConnection(&conf.DataBaseConfig)
	require.Nil(t, err)

	wagerRepo := repo.NewWagerRepo(conn)

	wagerService := services.NewWagerService(wagerRepo)

	wagerHandler := handlers.NewWagersHandler(wagerService)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wagerHandler.Handle)

	body, err := json.Marshal(placeWagerReq)
	require.Nil(t, err)

	req, err := http.NewRequest("POST", "/wagers", bytes.NewReader(body))
	require.Nil(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	var wager dto.Wager
	err = json.Unmarshal(rr.Body.Bytes(), &wager)
	require.Nil(t, err)
	require.NotNil(t, wager)
	require.NotEmpty(t, wager.ID)

	// 2. List wager
	req, err = http.NewRequest("GET", "/wagers?page=1&limit=20", nil)
	require.Nil(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	var wagerList []dto.Wager
	err = json.Unmarshal(rr.Body.Bytes(), &wagerList)

	require.Nil(t, err)
	require.NotEmpty(t, wagerList)

	idsMap := map[uint32]bool{}
	for _, w := range wagerList {
		idsMap[w.ID] = true
	}

	require.True(t, idsMap[wager.ID])

	// 3. Buy Wager
	buyWagerReq := dto.BuyWagerRequest{
		BuyingPrice: 20.5,
	}
	body, err = json.Marshal(buyWagerReq)
	require.Nil(t, err)

	purchaseRepo := repo.NewPurchaseRepo(conn)
	purchaseService := services.NewPurchaseService(purchaseRepo, wagerRepo)
	purchaseHandler := handlers.NewPurchasesHandler(purchaseService)
	purchaseHandle := http.HandlerFunc(purchaseHandler.Handle)

	req, err = http.NewRequest("POST", fmt.Sprintf("/buy/%d", wager.ID), bytes.NewReader(body))
	require.Nil(t, err)

	rr = httptest.NewRecorder()
	purchaseHandle.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
	var wagerPurchase dto.WagerPurchase
	err = json.Unmarshal(rr.Body.Bytes(), &wagerPurchase)
	require.Nil(t, err)
	require.NotNil(t, wagerPurchase)
	require.NotEmpty(t, wagerPurchase.ID)
	require.Equal(t, wager.ID, wagerPurchase.WagerID)
}

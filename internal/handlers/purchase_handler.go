package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/vitthalaa/wager-app/app_errors"
	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/services"
)

// PurchaseHandler is handler for all purchase(/buy) routes
type PurchaseHandler struct {
	purchaseService services.IPurchaseService
}

// NewPurchasesHandler ...
func NewPurchasesHandler(purchaseService services.IPurchaseService) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService: purchaseService,
	}
}

// Handle is method to handle requests to routes
func (h *PurchaseHandler) Handle(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case http.MethodPost:
		err = h.doPurchaseWager(w, req)
	default:
		log.Println("error no 404")
		writeResponse(w, http.StatusNotFound, app_errors.ErrorResponse{Code: app_errors.ErrNotFound})
	}

	if err != nil {
		log.Println("error {}", err)
		writeResponse(w, http.StatusInternalServerError, app_errors.ErrorResponse{Code: app_errors.ErrInternalError})
	}
}

func (h *PurchaseHandler) doPurchaseWager(w http.ResponseWriter, req *http.Request) error {
	id := strings.TrimPrefix(req.URL.Path, "/buy/")
	if id == "" {
		log.Println("empty wager id")
		writeResponse(w, http.StatusNotFound, app_errors.ErrorResponse{Code: app_errors.ErrNotFound})
		return nil
	}

	wagerID, err := strconv.Atoi(id)
	if err != nil {
		log.Println("invalid wager id {}", id)
		writeResponse(w, http.StatusNotFound, app_errors.ErrorResponse{Code: app_errors.ErrNotFound})
		return nil
	}

	decoder := json.NewDecoder(req.Body)
	var request dto.BuyWagerRequest
	err = decoder.Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &app_errors.ErrorResponse{Code: app_errors.ErrInvalidBody})
		return nil
	}

	request.WagerID = uint32(wagerID)

	res, err := h.purchaseService.PurchaseWager(req.Context(), &request)
	if err != nil {
		writeErrorResponse(w, err)
		return nil
	}

	writeResponse(w, http.StatusOK, res)
	return nil
}

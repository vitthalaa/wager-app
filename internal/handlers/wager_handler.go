package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/vitthalaa/wager-app/app_errors"
	"github.com/vitthalaa/wager-app/dto"
	"github.com/vitthalaa/wager-app/internal/services"
)

// WagersHandler is handler for all /wagers routes
type WagersHandler struct {
	wagerService services.IWagerService
}

// NewWagersHandler ...
func NewWagersHandler(wagerService services.IWagerService) *WagersHandler {
	return &WagersHandler{
		wagerService: wagerService,
	}
}

// Handle is method to handle requests to routes
func (h *WagersHandler) Handle(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case http.MethodPost:
		err = h.doPlaceWager(w, req)
	case http.MethodGet:
		err = h.doListWager(w, req)
	default:
		log.Println("error no 404")
		writeResponse(w, http.StatusNotFound, app_errors.ErrorResponse{Code: app_errors.ErrNotFound})
	}

	if err != nil {
		log.Println("error {}", err)
		writeResponse(w, http.StatusInternalServerError, app_errors.ErrorResponse{Code: app_errors.ErrInternalError})
	}
}

// doPlaceWager places wager
func (h *WagersHandler) doPlaceWager(w http.ResponseWriter, req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	var request dto.PlaceWagerRequest
	err := decoder.Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &app_errors.ErrorResponse{Code: app_errors.ErrInvalidBody})
		return nil
	}

	wager, err := h.wagerService.PlaceWager(req.Context(), &request)
	if err != nil {
		writeErrorResponse(w, err)
		return nil
	}

	writeResponse(w, http.StatusOK, wager)
	return nil
}

// doListWager list wager
func (h *WagersHandler) doListWager(w http.ResponseWriter, req *http.Request) error {
	pageStr := strings.TrimSpace(req.URL.Query().Get("page"))
	limitStr := strings.TrimSpace(req.URL.Query().Get("limit"))
	page, limit := 0, 0
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
	}

	if err != nil {
		log.Printf("get query params err %s", err)
	}

	request := &dto.ListWagerRequest{
		Page:  uint32(page),
		Limit: uint32(limit),
	}

	wagerList, err := h.wagerService.ListWager(req.Context(), request)
	if err != nil {
		writeErrorResponse(w, err)
		return nil
	}

	writeResponse(w, http.StatusOK, wagerList)
	return nil
}

func writeResponse(w http.ResponseWriter, status int, res interface{}) {
	resBody, err := json.Marshal(res)
	if err != nil {
		log.Println("marshal response body error {}", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(resBody)
	if err != nil {
		log.Println("write response body error {}", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal error")
	}
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*app_errors.ErrorResponse); ok {
		writeResponse(w, appErr.Status, appErr)
		return
	}

	writeResponse(w, http.StatusInternalServerError, app_errors.ErrorResponse{Code: app_errors.ErrInternalError})
}

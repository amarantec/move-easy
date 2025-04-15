package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/amarantec/move-easy/internal"
	"github.com/amarantec/move-easy/internal/bus"
)

type BusHandler struct {
	service bus.IBusService
}

func NewBusHandler(service bus.IBusService) *BusHandler {
	return &BusHandler{service: service}
}

func (h *BusHandler) InsertNewBusLine(w http.ResponseWriter, r *http.Request) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newBusLine := internal.BusLine{}
	if err :=
		json.NewDecoder(r.Body).Decode(&newBusLine); err != nil {
		http.Error(w,
			"could not decode this request, error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	response, err := h.service.InsertNewBusLine(ctxTimeout, newBusLine)
	if err != nil {
		http.Error(w,
			"could not insert this new line, error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

func (h *BusHandler) InsertBusStop(w http.ResponseWriter, r *http.Request) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newBusStop := internal.BusStop{}
	if err :=
		json.NewDecoder(r.Body).Decode(&newBusStop); err != nil {
		http.Error(w,
			"could not decode this request, error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	response, err := h.service.InsertBusStop(ctxTimeout, newBusStop)
	if err != nil {
		http.Error(w,
			"could not insert this bus stop, error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

func (h *BusHandler) GetBusLine(w http.ResponseWriter, r *http.Request) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	busLineID, err := strconv.ParseInt(r.PathValue("busLineID"), 10, 64)
	if err != nil {
		http.Error(w,
			"invalid parameter, error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	response, err := h.service.GetBusLine(ctxTimeout, busLineID)
	if err != nil {
		http.Error(w,
			"could not get this bus line, error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

func (h *BusHandler) GetBusStop(w http.ResponseWriter, r *http.Request) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	busStopID, err := strconv.ParseInt(r.PathValue("busStopID"), 10, 64)
	if err != nil {
		http.Error(w,
			"invalid parameter, error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	response, err := h.service.GetBusStop(ctxTimeout, busStopID)
	if err != nil {
		http.Error(w,
			"could not get this bus stop, error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

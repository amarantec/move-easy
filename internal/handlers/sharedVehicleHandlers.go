package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/amarantec/move-easy/internal"
	"github.com/amarantec/move-easy/internal/middleware"
	"github.com/amarantec/move-easy/internal/sharedVehicle"
)

type SharedVehicleHandler struct {
	service sharedVehicle.ISharedVehicleService
}

func NewSharedVehicleHandler(service sharedVehicle.ISharedVehicleService) *SharedVehicleHandler {
	return &SharedVehicleHandler{service: service}
}

func (h *SharedVehicleHandler) InsertSharedVehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,
			"invalid http method",
			http.StatusMethodNotAllowed)
		return
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := r.Context().Value(middleware.UserIDKey).(int64)

	sharedVehicle := internal.SharedVehicle{}
	sharedVehicle.UserID = userID

	if err :=
		json.NewDecoder(r.Body).Decode(&sharedVehicle); err != nil {
		http.Error(w,
			"coudl not decote this request, error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	response, err := h.service.InsertSharedVehicle(ctxTimeout, sharedVehicle)
	if err != nil {
		http.Error(w,
			"could not insert this vehicle, error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

func (h *SharedVehicleHandler) GetSharedVehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w,
			"invalid http method",
			http.StatusMethodNotAllowed)
		return
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vehicleID, err := strconv.ParseInt(r.PathValue("vehicleID"), 10, 64)
	if err != nil {
		http.Error(w,
			"invalid parameter, error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	response, err := h.service.GetSharedVehicle(ctxTimeout, vehicleID)
	if err != nil {
		http.Error(w,
			"could not get this vehicle, error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

func (h *SharedVehicleHandler) ListAllSharedVehicles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w,
			"invalid http method",
			http.StatusMethodNotAllowed)
		return
	}
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.service.ListAllSharedVehicles(ctxTimeout)
	if err != nil {
		http.Error(w,
			"could not list all available shared vehicles, error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

func (h *SharedVehicleHandler) UpdateSharedVehicleLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w,
			"invalid http method",
			http.StatusMethodNotAllowed)
		return
	}
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := r.Context().Value(middleware.UserIDKey).(int64)

	sharedVehicle := internal.SharedVehicle{}
	sharedVehicle.UserID = userID

	if err :=
		json.NewDecoder(r.Body).Decode(&sharedVehicle); err != nil {
		http.Error(w,
			"could not decode this request, error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateSharedVehicleLocation(ctxTimeout, sharedVehicle)
	if err != nil {
		http.Error(w,
			"could not update this shared vehicle location, error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

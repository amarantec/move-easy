package sharedVehicle

import (
	"context"
	"errors"

	"github.com/amarantec/move-easy/internal"
)

type ISharedVehicleService interface {
	InsertSharedVehicle(ctx context.Context, vehicle internal.SharedVehicle) (int64, error)
	ListAllSharedVehicles(ctx context.Context) ([]internal.SharedVehicle, error)
	GetSharedVehicle(ctx context.Context, vehicleID int64) (internal.SharedVehicle, error)
	UpdateSharedVehicleLocation(ctx context.Context, vehicle internal.SharedVehicle) (bool, error)
}

type sharedVehicleService struct {
	repository ISharedVehicleRepository
}

func NewSharedVehicleService(repo ISharedVehicleRepository) ISharedVehicleService {
	return &sharedVehicleService{repository: repo}
}

func (s *sharedVehicleService) InsertSharedVehicle(ctx context.Context, vehicle internal.SharedVehicle) (int64, error) {
	if valid, err := validateSharedVehicle(vehicle); err != nil || !valid {
		return internal.ZERO, err
	}
	return s.repository.InsertSharedVehicle(ctx, vehicle)
}

func (s *sharedVehicleService) ListAllSharedVehicles(ctx context.Context) ([]internal.SharedVehicle, error) {
	return s.repository.ListAllSharedVehicles(ctx)
}

func (s *sharedVehicleService) GetSharedVehicle(ctx context.Context, vehicleID int64) (internal.SharedVehicle, error) {
	return s.repository.GetSharedVehicle(ctx, vehicleID)
}

func (s *sharedVehicleService) UpdateSharedVehicleLocation(ctx context.Context, vehicle internal.SharedVehicle) (bool, error) {
	if valid, err := validateSharedVehicle(vehicle); err != nil || !valid {
		return false, err
	}
	return s.repository.UpdateSharedVehicleLocation(ctx, vehicle)
}

func validateSharedVehicle(sv internal.SharedVehicle) (bool, error) {
	if sv.UserID <= internal.ZERO {
		return false, ErrSVUserIDEmpty
	}
	if sv.VehicleType < internal.ZERO {
		return false, ErrSVTypeEmpty
	}

	return true, nil
}

var (
	ErrSVUserIDEmpty = errors.New("Erro shared vehicle user id empty")
	ErrSVTypeEmpty   = errors.New("Error shared vahicle type empty")
)

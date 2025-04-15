package bus

import (
	"context"

	"github.com/amarantec/move-easy/internal"
)

type IBusService interface {
	InsertNewBusLine(ctx context.Context, busline internal.BusLine) (int64, error)
	InsertBusStop(ctx context.Context, busStop internal.BusStop) (int64, error)
	GetBusLine(ctx context.Context, busLineID int64) (internal.BusLine, error)
	GetBusStop(ctx context.Context, busStopID int64) (internal.BusStop, error)
}

type busService struct {
	repository IBusRepository
}

func NewBusService(repo IBusRepository) IBusService {
	return &busService{repository: repo}
}

func (s *busService) InsertNewBusLine(ctx context.Context, busLine internal.BusLine) (int64, error) {
	return s.repository.InsertNewBusLine(ctx, busLine)
}

func (s *busService) InsertBusStop(ctx context.Context, busStop internal.BusStop) (int64, error) {
	return s.repository.InsertBusStop(ctx, busStop)
}

func (s *busService) GetBusLine(ctx context.Context, busLineID int64) (internal.BusLine, error) {
	return s.repository.GetBusLine(ctx, busLineID)
}

func (s *busService) GetBusStop(ctx context.Context, busStopID int64) (internal.BusStop, error) {
	return s.repository.GetBusStop(ctx, busStopID)
}

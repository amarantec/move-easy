package sharedVehicle

import (
	"context"

	"github.com/amarantec/move-easy/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ISharedVehicleRepository interface{}

type sharedVehicleRepository struct {
	Conn *pgxpool.Pool
}

func NewSharedVehicleRepository(connection *pgxpool.Pool) ISharedVehicle {
	return &sharedVehicleRepository{Conn: connection}
}

func (r *sharedVehicleRepository) InsertVehicle(ctx context.Context, vehicle internal.SharedVehicle) (int64, error) {
	if

}

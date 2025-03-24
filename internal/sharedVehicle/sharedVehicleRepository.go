package sharedVehicle

import (
	"context"
	"time"

	"github.com/amarantec/move-easy/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ISharedVehicleRepository interface {
	InsertSharedVehicle(ctx context.Context, vehicle internal.SharedVehicle) (int64, error)
	ListAllSharedVehicles(ctx context.Context) ([]internal.SharedVehicle, error)
	GetSharedVehicle(ctx context.Context, vehicleID int64) (internal.SharedVehicle, error)
	UpdateSharedVehicleLocation(ctx context.Context, vehicle internal.SharedVehicle) (bool, error)
}

type sharedVehicleRepository struct {
	Conn *pgxpool.Pool
}

func NewSharedVehicleRepository(connection *pgxpool.Pool) ISharedVehicleRepository {
	return &sharedVehicleRepository{Conn: connection}
}

func (r *sharedVehicleRepository) InsertSharedVehicle(ctx context.Context, vehicle internal.SharedVehicle) (int64, error) {
	if err :=
		r.Conn.QueryRow(
			ctx,
			`INSERT INTO shared_vehicle (user_id, latitude, longitude, vehicle_type, reported_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`, vehicle.UserID, vehicle.Latitude, vehicle.Longitude, vehicle.VehicleType, time.Now()).Scan(&vehicle.ID); err != nil {
		return internal.ZERO, err
	}
	return vehicle.ID, nil
}

func (r *sharedVehicleRepository) GetSharedVehicle(ctx context.Context, vehicleID int64) (internal.SharedVehicle, error) {
	vehicle := internal.SharedVehicle{}
	if err :=
		r.Conn.QueryRow(
			ctx,
			`SELECT latitude, longitude, vehicle_type, reported_at
				FROM shared_vehicle WHERE id = $1 AND deleted_at IS NULL;`, vehicleID).Scan(&vehicle.Latitude,
			&vehicle.Longitude, &vehicle.VehicleType, &vehicle.ReportedAt); err != nil {
		if err == pgx.ErrNoRows {
			return internal.SharedVehicle{}, nil
		}
		return internal.SharedVehicle{}, err
	}

	return vehicle, nil
}

func (r *sharedVehicleRepository) ListAllSharedVehicles(ctx context.Context) ([]internal.SharedVehicle, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	sharedVehicleChannel := make(chan []internal.SharedVehicle)
	errorChannel := make(chan error)

	go func() {
		rows, err :=
			r.Conn.Query(
				ctx,
				`SELECT id, latitude, longitude, vehicle_type, reported_at
				FROM shared_vehicle WHERE deleted_at IS NULL;`)
		if err != nil {
			errorChannel <- err
			return
		}

		defer rows.Close()
		vehicles := []internal.SharedVehicle{}
		for rows.Next() {
			v := internal.SharedVehicle{}
			if err := rows.Scan(
				&v.ID,
				&v.Latitude,
				&v.Longitude,
				&v.VehicleType,
				&v.ReportedAt); err != nil {
				errorChannel <- err
				return
			}
			vehicles = append(vehicles, v)
		}
		sharedVehicleChannel <- vehicles
	}()

	select {
	case vehicles := <-sharedVehicleChannel:
		return vehicles, nil
	case err := <-errorChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *sharedVehicleRepository) UpdateSharedVehicleLocation(ctx context.Context, vehicle internal.SharedVehicle) (bool, error) {
	result, err :=
		r.Conn.Exec(
			ctx,
			`UPDATE shared_vehicle SET user_id = $2, latitude = $3, longitude = $4, reported_at = $5, updated_at = $6
			WHERE id = $1 AND deleted_at IS NULL;`, vehicle.ID, vehicle.UserID, vehicle.Latitude, vehicle.Longitude, time.Now(), time.Now())

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == internal.ZERO {
		return false, nil
	} else {
		return true, nil
	}
}

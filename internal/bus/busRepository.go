package bus

import (
	"context"

	"github.com/amarantec/move-easy/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IBusRepository interface {
	InsertNewBusLine(ctx context.Context, busline internal.BusLine) (int64, error)
	InsertBusStop(ctx context.Context, busStop internal.BusStop) (int64, error)
	GetBusLine(ctx context.Context, busLineID int64) (internal.BusLine, error)
	GetBusStop(ctx context.Context, busStopID int64) (internal.BusStop, error)
}

type busRepository struct {
	Conn *pgxpool.Pool
}

func NewBusRepository(connection *pgxpool.Pool) IBusRepository {
	return &busRepository{Conn: connection}
}

func (r *busRepository) InsertNewBusLine(ctx context.Context, busline internal.BusLine) (int64, error) {
	if err :=
		r.Conn.QueryRow(
			ctx,
			`INSERT INTO bus_line (name, bus_init, bus_end) VALUES ($1, $2, $3) 
				RETURNING id;`, busline.Name, busline.BusInit.ID, busline.BusEnd.ID).Scan(&busline.ID); err != nil {
		return internal.ZERO, err
	}

	return busline.ID, nil
}

func (r *busRepository) GetBusLine(ctx context.Context, busLineID int64) (internal.BusLine, error) {
	busLine := internal.BusLine{ID: busLineID}

	if err :=
		r.Conn.QueryRow(ctx,
			`SELECT name, bus_init, bus_end WHERE id = $1
				AND deleted_at IS NULL;`, busLineID).Scan(&busLine.Name, busLine.BusInit.ID,
			&busLine.BusEnd.ID); err != nil {
		return internal.BusLine{}, err
	}
	// BUSCAR DETALHES DOS PONTOS DE PARADA NA ROTA
	busLine.BusInit, _ = r.GetBusStop(ctx, busLine.BusInit.ID)
	busLine.BusEnd, _ = r.GetBusStop(ctx, busLine.BusEnd.ID)

	if busLine.BusInit.Name == internal.EMPTY || busLine.BusEnd.Name == internal.EMPTY {
		return internal.BusLine{}, nil
	}

	return busLine, nil
}

func (r *busRepository) InsertBusStop(ctx context.Context, busStop internal.BusStop) (int64, error) {
	if err :=
		r.Conn.QueryRow(
			ctx,
			`INSERT INTO bus_stop (name, latitude, longitude) VALUES ($1, $2, $3) 
				RETURNING id;`, busStop.Name, busStop.Latitude, busStop.Longitude).Scan(&busStop.ID); err != nil {
		return internal.ZERO, err
	}

	return busStop.ID, nil
}

func (r *busRepository) GetBusStop(ctx context.Context, busStopID int64) (internal.BusStop, error) {
	busStop := internal.BusStop{ID: busStopID}
	if err :=
		r.Conn.QueryRow(
			ctx,
			`SELECT name, latitude, longitude FROM bus_stop WHERE id= $1
				AND deleted_at IS NULL;`, busStopID).Scan(&busStop.Name,
			&busStop.Latitude, &busStop.Longitude); err != nil {
		return internal.BusStop{}, err
	}
	return busStop, nil
}

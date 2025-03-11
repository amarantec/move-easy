package address

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/amarantec/move-easy/internal"
	"log"
)

type IAddressRepository interface {
    AddOrUpdateAddress(ctx context.Context, address internal.Address) (int64, error)
    GetAddress(ctx context.Context, userID int64) (internal.Address, error)
}

type addressRepository struct {
	Conn	*pgxpool.Pool
}

func NewAddressRepository(conn *pgxpool.Pool) IAddressRepository {
	return &addressRepository{Conn: conn}
}

func (r *addressRepository) GetAddress(ctx context.Context, userID int64) (internal.Address, error) {
	address := internal.Address{UserID: userID}
	if err :=
		r.Conn.QueryRow(
			ctx,
			`SELECT id, street, number, cep, neighborhood, city, state
				FROM address WHERE user_id = $1;`, userID).Scan(&address.ID,
					&address.Street, &address.Number, &address.CEP,
					&address.Neighborhood, &address.City, &address.State); err != nil {

		if err == pgx.ErrNoRows {
			return internal.Address{}, nil
		}

		return internal.Address{}, err
	}

	return address, nil
}

func (r *addressRepository) AddOrUpdateAddress(ctx context.Context, address internal.Address) (int64, error) {
	addressDB, err := r.GetAddress(ctx, address.UserID)
	if err != nil && err != pgx.ErrNoRows {
		return internal.ZERO, err
	}

	if addressDB.ID == internal.ZERO {
		err := r.Conn.QueryRow(
			ctx,
			`INSERT INTO address (user_id, street, number, cep, neighborhood, city, state) VALUES
				($1, $2, $3, $4, $5, $6, $7) RETURNING id;`, address.UserID, address.Street, address.Number,
					address.CEP, address.Neighborhood, address.City, address.State).Scan(&address.ID)
		if err != nil {
			return internal.ZERO, err
		}
		return address.ID, nil
	} else {
		res, err :=
			r.Conn.Exec(
				ctx,
				`UPDATE address SET street = $3,
					number = $4,
					cep = $5,
					neighborhood = $6,
					city = $7,
					state = $8
				WHERE id = $1 AND user_id = $2;`, address.ID, address.UserID, address.Street, address.Number,
					address.CEP, address.Neighborhood, address.City, address.State)
				if err != nil {
					return internal.ZERO, err
				}

				if res.RowsAffected() == internal.ZERO {
					log.Printf("%d rows affected", res.RowsAffected())
					return internal.ZERO, nil
				} else {
					log.Printf("%d rows affected", res.RowsAffected())
					return address.ID, nil
				}
	}
}



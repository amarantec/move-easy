package db

import "context"

func createTables(ctx context.Context) {
	createUsersTable := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            first_name VARCHAR(100),
            last_name VARCHAR(100),
            email VARCHAR(100) UNIQUE,
            password TEXT,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP NULL,
            deleted_at TIMESTAMP NULL
        );`

	_, err := Conn.Exec(ctx, createUsersTable)
	if err != nil {
		panic(err)
	}

	createAddressTable := `
        CREATE TABLE IF NOT EXISTS address (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users (id),
            street VARCHAR(100),
            number TEXT,
            cep VARCHAR(8),
            neighborhood VARCHAR(100),
            city VARCHAR(100),
            state VARCHAR(2),
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP NULL,
            deleted_at TIMESTAMP NULL
        );`

	_, err = Conn.Exec(ctx, createAddressTable)
	if err != nil {
		panic(err)
	}

	createContactTable := `
        CREATE TABLE IF NOT EXISTS contacts (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users (id),
            name    VARCHAR(100),
            ddi     VARCHAR(3),
            ddd     VARCHAR(3),
            phone_number VARCHAR(9),
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP NULL,
            deleted_at TIMESTAMP NULL
        );`
	_, err = Conn.Exec(ctx, createContactTable)
	if err != nil {
		panic(err)
	}

	createSharedVehicleTable := `
		CREATE TABLE IF NOT EXISTS shared_vehicle (
			id SERIAL PRIMARY KEY,
			user_id	INTEGER REFERENCES users (id),
			latitude DOUBLE PRECISION,
			longitude DOUBLE PRECISION,
			vehicle_type INTEGER,
			reported_at TIMESTAMP DEFAULT NOW(),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP NULL,
			deleted_at TIMESTAMP NULL
		);`

	_, err = Conn.Exec(ctx, createSharedVehicleTable)
	if err != nil {
		panic(err)
	}
}

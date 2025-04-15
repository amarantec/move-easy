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

	createBusStopTable := `
		CREATE TABLE IF NOT EXISTS bus_stop (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP NULL,
			deleted_at TIMESTAMP NULL
		);`

	_, err = Conn.Exec(ctx, createBusStopTable)
	if err != nil {
		panic(err)
	}

	createBusLineTable := `
		CREATE TABLE IF NOT EXISTS bus_line (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			bus_init INTEGER NOT NULL REFERENCES bus_stop(id),
			bus_end  INTEGER NOT NULL REFERENCES bus_stop(id),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP NULL,
			deleted_at TIMESTAMP NULL
		);`

	_, err = Conn.Exec(ctx, createBusLineTable)
	if err != nil {
		panic(err)
	}

	createBusScheduleTable := `
		CREATE TABLE IF NOT EXISTS bus_schedule (
			id SERIAL PRIMARY KEY,
			bus_line_id INTEGER NOT NULL REFERENCES bus_line(id),
			day_of_week VARCHAR(20) NOT NULL,
			start_time 	TIME NOT NULL,
			end_time    TIME NOT NULL,
			created_at	TIMESTAMP DEFAULT NOW(),
			updated_at 	TIMESTAMP NULL,
			deleted_at  TIMESTAMP NULL
		);`

	_, err = Conn.Exec(ctx, createBusScheduleTable)
	if err != nil {
		panic(err)
	}

	createMetroTable := `
		CREATE TABLE IF NOT EXISTS metro (
			id SERIAL PRIMARY KEY,
			station_name VARCHAR(255) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP NULL,
			deleted_at TIMESTAMP NULL
		);`

	_, err = Conn.Exec(ctx, createMetroTable)
	if err != nil {
		panic(err)
	}
}

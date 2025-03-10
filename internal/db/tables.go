package db

import "context"

func createTables(ctx context.Context) {
    createUsersTable :=`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            first_name VARCHAR(100),
            last_name VARCHAR(100),
            email VARCHAR(100) UNIQUE,
            password TEXT
        );`

    _, err := Conn.Exec(ctx, createUsersTable)
    if err != nil {
        panic(err)
    }

    createAddressTable :=`
        CREATE TABLE IF NOT EXISTS address (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users (id),
            street VARCHAR(100),
            number TEXT,
            cep VARCHAR(8),
            neighborhood VARCHAR(100),
            city VARCHAR(100),
            state VARCHAR(2)
        );`

    _, err = Conn.Exec(ctx, createAddressTable)
    if err != nil {
        panic(err)
    }
}


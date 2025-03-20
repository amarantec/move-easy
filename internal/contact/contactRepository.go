package contact

import (
    "context"
    "time"
    "github.com/amarantec/move-easy/internal"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
)

type IContactRepository interface {
    SaveContact(ctx context.Context, contact internal.Contact) (int64, error)
    GetContact(ctx context.Context, userID, contactID int64) (internal.Contact, error)
    ListContacts(ctx context.Context, userID int64) ([]internal.Contact, error)
    UpdateContact(ctx context.Context, contact internal.Contact) (bool, error)
    DeleteContact(ctx context.Context, userID, contactID int64) (bool, error)
}

type contactRepository struct {
    Conn    *pgxpool.Pool
}

func NewContactRepository(connection *pgxpool.Pool) IContactRepository {
    return &contactRepository{Conn: connection}
}

func (r *contactRepository) SaveContact(ctx context.Context, contact internal.Contact) (int64, error) {
    err :=
        r.Conn.QueryRow(
            ctx,
            `INSERT INTO contacts (user_id, name, ddi, ddd, phone_number) VALUES
                ($1, $2, $3, $4, $5) RETURNING id;`, contact.UserID, contact.Name, contact.DDI,
                    contact.DDD, contact.PhoneNumber).Scan(&contact.ID)
    if err != nil {
        return internal.ZERO, err
    }

    return contact.ID, nil
}

func (r *contactRepository) GetContact(ctx context.Context, userID, contactID int64) (internal.Contact, error) {
    contact := internal.Contact{ID: contactID, UserID: userID}
    if err :=
        r.Conn.QueryRow(
            ctx,
            `SELECT name, ddi, ddd, phone_number FROM contacts
                WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL;`, userID, contactID).Scan(&contact.Name,
                    &contact.DDI, &contact.DDD, &contact.PhoneNumber); err != nil {

        if err == pgx.ErrNoRows {
            return internal.Contact{}, nil
        }

        return internal.Contact{}, err
    }

    return contact, nil
}

func (r *contactRepository) ListContacts(ctx context.Context, userID int64) ([]internal.Contact, error) {
    rows, err :=
        r.Conn.Query(
            ctx,
            `SELECT id, name, ddi, ddd, phone_number FROM contacts
                WHERE user_id = $1 AND deleted_at IS NULL`, userID)
    if err != nil {
        return []internal.Contact{}, err
    }
    defer rows.Close()

    contacts := []internal.Contact{}
    for rows.Next() {
        c := internal.Contact{}
        if err := rows.Scan(
            &c.ID,
            &c.Name,
            &c.DDI,
            &c.DDD,
            &c.PhoneNumber); err != nil {
                return []internal.Contact{}, err
        }
        c.UserID = userID
        contacts = append(contacts, c)
    }

    if err := rows.Err(); err != nil {
        return []internal.Contact{}, err
    }

    return contacts, nil
}

func (r *contactRepository) UpdateContact(ctx context.Context, contact internal.Contact) (bool, error) {
    result, err :=
        r.Conn.Exec(
            ctx,
            `UPDATE contacts SET name = $3, ddi = $4, ddd = $5,
                phone_number = $6, updated_at = $7 WHERE user_id = $1 AND id = $2
                    AND deleted_at IS NULL;`, contact.UserID, contact.ID, contact.Name, contact.DDI, contact.DDD, contact.PhoneNumber, time.Now())
    if err != nil {
        return false, err
    }

    if result.RowsAffected() == internal.ZERO {
        return false, nil
    } else {
        return true, nil
    }
}

func (r *contactRepository) DeleteContact(ctx context.Context, userID, contactID int64) (bool, error) {
    result, err :=
        r.Conn.Exec(
            ctx,
            `UPDATE contacts SET deleted_at = $3 WHERE user_id = $1 
                AND id = $2 AND deleted_at IS NULL;`, userID, contactID,
                time.Now())
    if err != nil {
        return false, err
    }

    if result.RowsAffected() == internal.ZERO {
        return false, nil
    } else {
        return true, nil
    }
}

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/LSariol/Cove/internal/crypt"
)

// CreateSecret takes a secret and stores it into the Cove database
func (d *Database) CreateSecret(ctx context.Context, s Secret) (Secret, error) {
	var newSecret Secret

	const query = `
	INSERT INTO cove.secrets (secret_key, secret_value)
	Values ($1, $2)
	RETURNING secret_key, date_added;
	`

	encryptedKey, err := crypt.Encrypt(s.Value)
	if err != nil {
		return newSecret, fmt.Errorf("encrypt: %w", err)
	}
	s.Value = encryptedKey

	row := d.Pool.QueryRow(ctx, query, s.Key, s.Value)

	err = row.Scan(
		&newSecret.Key,
		&newSecret.DateAdded)

	if err != nil {
		return newSecret, fmt.Errorf("row.scan: %w", err)
	}

	fmt.Printf("Inserted secret with key %s at %s\n", newSecret.Key, newSecret.DateAdded.Format(time.RFC3339))
	return newSecret, nil
}

func (d *Database) GetSecret(ctx context.Context, key string) (Secret, error) {
	var s Secret

	tx, _ := d.Pool.Begin(ctx)
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE cove.secrets
		SET 
			times_pulled = times_pulled + 1
		WHERE secret_key = $1
	`, key); err != nil {
		return s, err
	}

	row := d.Pool.QueryRow(ctx, `
	SELECT id, secret_key, secret_value, version, times_pulled, date_added, last_modified
	FROM cove.secrets 
	WHERE secret_key=$1
	`, key)

	err := row.Scan(
		&s.Id,
		&s.Key,
		&s.Value,
		&s.Version,
		&s.TimesPulled,
		&s.DateAdded,
		&s.LastModified)

	if err != nil {
		return s, fmt.Errorf("row.scan %q: %w", key, err)
	}

	decryptedVal, err := crypt.Decrypt(s.Value)
	if err != nil {
		return s, fmt.Errorf("get secret %q: %w", key, err)
	}

	s.Value = decryptedVal

	return s, tx.Commit(ctx)
}

func (d *Database) GetAllKeys(ctx context.Context) ([]Secret, error) {
	var secrets []Secret

	const query = `
	SELECT id, secret_key, secret_value, version, times_pulled, date_added, last_modified
	FROM cove.secrets 
	ORDER BY secret_key ASC
	`

	rows, err := d.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query secrets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s Secret
		if err := rows.Scan(
			&s.Id,
			&s.Key,
			&s.Value,
			&s.Version,
			&s.TimesPulled,
			&s.DateAdded,
			&s.LastModified,
		); err != nil {
			return nil, fmt.Errorf("scan secret: %w", err)
		}

		s.Value = ""
		secrets = append(secrets, s)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return secrets, nil
}

func (d *Database) UpdateSecret(ctx context.Context, s Secret) error {
	const query = `
	UPDATE cove.secrets
	SET
		secret_value = $2,
		version = version + 1
	WHERE secret_key = $1
	`
	encryptedVal, err := crypt.Encrypt(s.Value)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	cmdTag, err := d.Pool.Exec(ctx, query, s.Key, encryptedVal)
	if err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	fmt.Println("Rows affected: ", cmdTag.RowsAffected())
	return nil
}

func (d *Database) DeleteSecret(ctx context.Context, key string) error {
	const query = `
	DELETE FROM cove.secrets
	WHERE secret_key = $1;
	`

	cmdTag, err := d.Pool.Exec(ctx, query, key)
	if err != nil {
		return fmt.Errorf("delete secret %q: %w", key, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no secret found with key %q", key)
	}
	return nil

}

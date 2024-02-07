package website

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Migrate(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS websites (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL,
		rank INT NOT NULL
	)
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, website Website) (*Website, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, "INSERT INTO websites(name, url, rank) values ($1, $2, $3) RETURNING id", website.Name, website.URL, website.Rank).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	website.ID = id

	return &website, nil
}

func (r *PostgresRepository) All(ctx context.Context) ([]Website, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM websites")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var websites []Website
	for rows.Next() {
		var website Website
		if err := rows.Scan(&website.ID, &website.Name, &website.URL, &website.Rank); err != nil {
			return nil, err
		}

		websites = append(websites, website)
	}

	return websites, nil
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (*Website, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM websites WHERE name = $1", name)

	var website Website
	if err := row.Scan(&website.ID, &website.Name, &website.URL, &website.Rank); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &website, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id int64, updatedWebsite Website) (*Website, error) {
	result, err := r.db.ExecContext(ctx, "UPDATE websites SET name = $1, url = $2, rank = $3 WHERE id = $4", updatedWebsite.Name, updatedWebsite.URL, updatedWebsite.Rank, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdatedFailed
	}

	return &updatedWebsite, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM websites WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeletedFailed
	}

	return err
}

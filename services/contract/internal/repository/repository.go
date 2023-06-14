package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"microservices/services/contract/internal/domain"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Repo struct {
	db *pgxpool.Pool
}

type Contract interface {
	Create(ctx context.Context, contract *domain.Contract) error
	GetByID(ctx context.Context, id int64) (*domain.Contract, error)
	GetAll(ctx context.Context, title string, filters Filters) ([]*domain.Contract, error)
	Delete(ctx context.Context, id int64) error
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (s *Repo) Create(ctx context.Context, contract *domain.Contract) error {
	query := `
		INSERT INTO contracts (title, description)
		VALUES ($1, $2)
		RETURNING id, created_at, version`

	args := []interface{}{contract.Title, contract.Desc}

	return s.db.QueryRow(ctx, query, args...).Scan(&contract.ID, &contract.CreatedAt, &contract.Version)
}

func (s *Repo) GetByID(ctx context.Context, id int64) (*domain.Contract, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, description, version
		FROM contracts
		WHERE id = $1`

	var contract domain.Contract

	err := s.db.QueryRow(ctx, query, id).Scan(
		&contract.ID,
		&contract.CreatedAt,
		&contract.Title,
		&contract.Desc,
		&contract.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &contract, nil
}

func (s *Repo) GetAll(ctx context.Context, title string, filters Filters) ([]*domain.Contract, error) {
	query := fmt.Sprintf(`
		SELECT id, created_at, title, description, version
		FROM contracts
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	args := []any{title, filters.limit(), filters.offset()}

	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	defer rows.Close()

	contracts := []*domain.Contract{}

	for rows.Next() {
		var contract domain.Contract

		err := rows.Scan(
			&contract.ID,
			&contract.CreatedAt,
			&contract.Title,
			&contract.Desc,
			&contract.Version,
		)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, &contract)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contracts, nil
}

func (s *Repo) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM contracts
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

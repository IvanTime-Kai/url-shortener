package repository

import (
	"context"

	"github.com/IvanTime-Kai/url-shortener/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (r *LinkRepository) Create(ctx context.Context, link *domain.Link) error {
	query := `
		INSERT INTO links (original_url, code, created_at, expires_at)
		VALUES ($1, $2, NOW(), $3)
		RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query, link.OriginalURL, link.Code, link.ExpiresAt).Scan(&link.ID, &link.CreatedAt)
}

func (r *LinkRepository) FindByCode(ctx context.Context, code string) (*domain.Link, error) {
	query := `SELECT id, original_url, code, created_at, expires_at FROM links WHERE code = $1`

	link := &domain.Link{}

	err := r.db.QueryRow(ctx, query, code).Scan(&link.ID, &link.OriginalURL, &link.Code, &link.CreatedAt, &link.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (r *LinkRepository) FindAll(ctx context.Context) ([]domain.Link, error) {
	query := `SELECT id, original_url, code, created_at, expires_at FROM links ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []domain.Link

	for rows.Next() {
		var link domain.Link

		if err := rows.Scan(&link.ID, &link.OriginalURL, &link.Code, &link.CreatedAt, &link.ExpiresAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (r *LinkRepository) DeleteByCode(ctx context.Context, code string) error {
	query := `DELETE FROM links WHERE code = $1`

	_, err := r.db.Exec(ctx, query, code)
	return err 
}
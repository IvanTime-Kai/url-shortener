package repository

import (
	"context"

	"github.com/IvanTime-Kai/url-shortener/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClickRepository struct {
	db *pgxpool.Pool
}

func NewClickRepository(db *pgxpool.Pool) *ClickRepository {
	return &ClickRepository{
		db: db,
	}
}

func (r *ClickRepository) Create(ctx context.Context, click *domain.Click) error {
	query := `
		INSERT INTO clicks (link_id, ip, user_agent, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query, click.LinkID, click.IP, click.UserAgent).Scan(&click.ID, &click.CreatedAt)
}

func (r *ClickRepository) CountByLinkID(ctx context.Context, linkID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM clicks WHERE link_id = $1`
	err := r.db.QueryRow(ctx, query, linkID).Scan(&count)
	return count, err
}

func (r *ClickRepository) FindByLinkID(ctx context.Context, linkID int64) ([]domain.Click, error) {
	query := `
		SELECT id, link_id, ip, user_agent, created_at
		FROM clicks
		WHERE link_id = $1
		ORDER BY created_at desc
		LIMIT 100
	`

	rows, err := r.db.Query(ctx, query, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []domain.Click
	for rows.Next() {
		var click domain.Click
		if err := rows.Scan(&click.ID, &click.LinkID, &click.IP, &click.UserAgent, &click.CreatedAt); err != nil {
			return nil, err
		}

		clicks = append(clicks, click)
	}

	return clicks, nil
}

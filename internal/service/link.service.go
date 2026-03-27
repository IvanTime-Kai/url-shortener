package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/IvanTime-Kai/url-shortener/internal/domain"
	"github.com/IvanTime-Kai/url-shortener/internal/repository"
)

type LinkService struct {
	repo    *repository.LinkRepository
	baseURL string
}

func NewLinkService(repo *repository.LinkRepository, baseURL string) *LinkService {
	return &LinkService{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (s *LinkService) Shorten(ctx context.Context, originURL string) (*domain.Link, error) {
	if _, err := url.ParseRequestURI(originURL); err != nil {
		return nil, err
	}

	// generate code
	code, err := generateCode(6)
	if err != nil {
		return nil, err
	}

	link := &domain.Link{
		Code:        code,
		OriginalURL: originURL,
	}

	if err := s.repo.Create(ctx, link); err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	return link, nil
}

func (s *LinkService) Resolve(ctx context.Context, code string) (*domain.Link, error) {
	link, err := s.repo.FindByCode(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("link not found: %w", err)
	}

	return link, nil
}

func (s *LinkService) List(ctx context.Context) ([]domain.Link, error) {
	return s.repo.FindAll(ctx)
}

func (s *LinkService) Delete(ctx context.Context, code string) error {
	return s.repo.Delete(ctx, code)
}

func generateCode(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

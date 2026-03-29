package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/IvanTime-Kai/url-shortener/internal/cache"
	"github.com/IvanTime-Kai/url-shortener/internal/domain"
	"github.com/IvanTime-Kai/url-shortener/internal/repository"
)

type LinkService struct {
	repo      *repository.LinkRepository
	clickRepo *repository.ClickRepository

	counter *ClickCounter
	cache   *cache.LinkCache

	baseURL string
}

func NewLinkService(repo *repository.LinkRepository, clickRepo *repository.ClickRepository, counter *ClickCounter, cache *cache.LinkCache, baseURL string) *LinkService {
	return &LinkService{
		repo:      repo,
		clickRepo: clickRepo,
		baseURL:   baseURL,
		counter:   counter,
		cache:     cache,
	}
}

func (s *LinkService) Shorten(ctx context.Context, originURL string, ttlDays int) (*domain.Link, error) {
	if _, err := url.ParseRequestURI(originURL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// generate code
	code, err := generateCode(6)
	if err != nil {
		return nil, fmt.Errorf("generate code: %w", err)
	}

	link := &domain.Link{
		Code:        code,
		OriginalURL: originURL,
	}

	// set TTL if have
	if ttlDays > 0 {
		expires := time.Now().Add(time.Duration(ttlDays) * 24 * time.Hour)
		link.ExpiresAt = &expires
	}

	if err := s.repo.Create(ctx, link); err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	return link, nil
}

func (s *LinkService) Resolve(ctx context.Context, code string, ip string, userAgent string) (*domain.Link, error) {

	// Find in cache
	if link, err := s.cache.Get(ctx, code); err == nil {
		// check link expired
		if link.IsExpired() {
			s.cache.Delete(ctx, code)
			return nil, fmt.Errorf("link expired")
		}
		go s.recordClick(link, ip, userAgent)
		return link, nil
	}

	// cache miss => find DB
	link, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("link not found: %w", err)
	}

	// Check expired
	if link.IsExpired() {
		return nil, fmt.Errorf("link expired")
	}

	// save into cache => ready for next time
	go func() {
		s.cache.Set(ctx, link)
		s.recordClick(link, ip, userAgent)
	}()

	return link, nil
}

func (s *LinkService) GetStats(ctx context.Context, code string) (*domain.LinkStats, error) {
	link, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("link not found: %w", err)
	}

	memCount := s.counter.GetCount(link.ID)
	dbCount, err := s.clickRepo.CountByLinkID(ctx, link.ID)

	if err != nil {
		return nil, err
	}

	clicks, err := s.clickRepo.FindByLinkID(ctx, link.ID)
	if err != nil {
		return nil, err
	}

	return &domain.LinkStats{
		Link:        *link,
		TotalClicks: dbCount + memCount,
		Clicks:      clicks,
	}, nil
}

func (s *LinkService) List(ctx context.Context) ([]domain.Link, error) {
	return s.repo.FindAll(ctx)
}

func (s *LinkService) Delete(ctx context.Context, code string) error {
	if err := s.repo.DeleteByCode(ctx, code); err != nil {
		return err
	}
	return s.cache.Delete(ctx, code)
}

// =================== Helpers =======================

func generateCode(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

func (s *LinkService) recordClick(link *domain.Link, ip string, userAgent string) {
	click := &domain.Click{
		LinkID:    link.ID,
		IP:        ip,
		UserAgent: userAgent,
	}

	s.clickRepo.Create(context.Background(), click)
	s.counter.Increment(click.LinkID)
}

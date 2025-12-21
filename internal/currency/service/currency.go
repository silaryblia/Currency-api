package service

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"time"
)

type CurrencyService struct {
	repo     domain.CurrencyRepository
	provider domain.RatesProvider
}

func NewCurrencyService(repo domain.CurrencyRepository, provider domain.RatesProvider) *CurrencyService {
	return &CurrencyService{
		repo:     repo,
		provider: provider}
}

func (s *CurrencyService) GetAll() (map[string]domain.Currency, error) {
	return s.repo.GetAll()
}

func (s *CurrencyService) GetOne(code string) (domain.Currency, error) {
	return s.repo.GetOne(code)
}

func (s *CurrencyService) Create(code string, rate float64, date time.Time) error {
	return s.repo.Create(code, rate, date)
}

func (s *CurrencyService) UpdateOne(code string, rate float64, date time.Time) error {
	return s.repo.UpdateOne(code, rate, date)
}

func (s *CurrencyService) UpdateAll() error {
	s.repo.UpdateAll()
	return nil
}

func (s *CurrencyService) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s *CurrencyService) SyncRates(ctx context.Context) error {
	rates, rateDate, err := s.provider.ForceRefresh(ctx)
	if err != nil {
		return err
	}

	for code, rate := range rates {
		if err := s.repo.Upsert(code, rate, rateDate); err != nil {
			return err
		}
	}

	return nil
}

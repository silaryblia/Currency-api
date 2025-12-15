package service

import "Currency-apiNew2/internal/currency/domain"

type CurrencyService struct {
	repo domain.CurrencyRepository
}

func NewCurrencyService(repo domain.CurrencyRepository) *CurrencyService {
	return &CurrencyService{repo: repo}
}

func (s *CurrencyService) GetAll() (map[string]float64, error) {
	//return s.repo.GetAll()
	return s.GetAll()
}

func (s *CurrencyService) GetOne(code string) (rate float64, err error) {
	return s.repo.GetOne(code)
}

func (s *CurrencyService) Create(code string, rate float64) error {
	return s.repo.Create(code, rate)
}

func (s *CurrencyService) UpdateOne(code string, rate float64) error {
	return s.repo.UpdateOne(code, rate)
}

func (s *CurrencyService) UpdateAll() error {
	s.repo.UpdateAll()
	return nil
}

func (s *CurrencyService) DeleteAll() error {
	return s.repo.DeleteAll()
}

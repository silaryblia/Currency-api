package repository

import (
	"Currency-apiNew2/internal/currency/domain"
	"math/rand"
	"strings"
	"sync"

	"go.uber.org/zap"
)

const (
	DefaultUSD = 80.00
	DefaultEUR = 85.00
	DefaultAED = 20.00
)

type CurrencyRepoInMemory struct {
	mu     sync.RWMutex
	data   map[string]float64
	logger *zap.Logger
}

func NewCurrencyRepoInMemory(logger *zap.Logger) *CurrencyRepoInMemory {
	return &CurrencyRepoInMemory{
		data: map[string]float64{
			"USD": DefaultUSD,
			"EUR": DefaultEUR,
			"AED": DefaultAED,
		},
		logger: logger,
	}
}

func (r *CurrencyRepoInMemory) GetOne(code string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	v, ok := r.data[code]
	if !ok {
		return 0, domain.ErrNotFound
	}

	r.logger.Debug("repo: Get success", zap.String("code", code), zap.Float64("rate", v))
	return v, nil
}

func (r *CurrencyRepoInMemory) GetAll() (map[string]float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copy := make(map[string]float64, len(r.data))
	for k, v := range r.data {
		copy[k] = v
	}

	return copy, nil
}

// add new currency
func (r *CurrencyRepoInMemory) Create(code string, rate float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	if _, ok := r.data[code]; ok {
		return domain.ErrAlreadyExists
	}

	r.data[code] = rate
	return nil
}

// update one currency
func (r *CurrencyRepoInMemory) UpdateOne(code string, rate float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	if _, ok := r.data[code]; !ok {
		return domain.ErrNotFound
	}

	r.data[code] = rate
	return nil

}

// Update
func (r *CurrencyRepoInMemory) UpdateAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// clear old, add new
	for k, v := range r.data {
		change := rand.Float64()*10 - 5 // [-5, 5]
		newVal := v + change
		if newVal < 0 {
			newVal = 0
		}
		r.data[k] = newVal
	}
	return nil
}

// Delete
func (r *CurrencyRepoInMemory) DeleteAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// clear map
	r.data = make(map[string]float64)
	return nil
}

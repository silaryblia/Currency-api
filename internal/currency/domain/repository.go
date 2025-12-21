package domain

import "time"

type CurrencyRepository interface {
	GetOne(code string) (Currency, error)
	GetAll() (map[string]Currency, error)
	Create(code string, rate float64, date time.Time) error
	UpdateOne(code string, rate float64, date time.Time) error
	UpdateAll() error
	DeleteAll() error

	Upsert(code string, rate float64, date time.Time) error
}

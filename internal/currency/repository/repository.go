package repository

type CurrencyRepository interface {
	Get(code string) (float64, error)
	GetAll() (map[string]float64, error)
	Create(code string, rate float64) error
	UpdateOne(code string, rate float64) error
	UpdateAll() error
	DeleteAll() error
}

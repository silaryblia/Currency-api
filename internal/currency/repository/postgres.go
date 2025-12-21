package repository

import (
	"Currency-apiNew2/internal/currency/domain"
	"database/sql"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type CurrencyRepoPostgres struct {
	mu     sync.RWMutex
	db     *sql.DB
	logger *zap.Logger
}

func NewCurrencyRepoPostgres(db *sql.DB, logger *zap.Logger) *CurrencyRepoPostgres {
	return &CurrencyRepoPostgres{db: db, logger: logger}
}

func (r *CurrencyRepoPostgres) Upsert(
	code string,
	rate float64,
	rateDate time.Time,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	_, err := r.db.Exec(`
		INSERT INTO currencies (code, rate, rate_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (code) DO UPDATE SET
		    rate = EXCLUDED.rate,
		    rate_date = EXCLUDED.rate_date
		
	`,
		code, rate, rateDate)

	return err
}

func (r *CurrencyRepoPostgres) GetOne(code string) (domain.Currency, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	var c domain.Currency

	err := r.db.QueryRow(`
		SELECT code, rate, rate_date
		FROM currencies
		WHERE code = $1
	`, code).Scan(&c.Code, &c.Rate, &c.RateDate)

	if err == sql.ErrNoRows {
		return c, domain.ErrNotFound
	}

	return c, err
}

func (r *CurrencyRepoPostgres) GetAll() (map[string]domain.Currency, error) {
	rows, err := r.db.Query(`
		SELECT code, rate, rate_date
		FROM currencies
		ORDER BY code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]domain.Currency)

	for rows.Next() {
		var c domain.Currency
		if err := rows.Scan(&c.Code, &c.Rate, &c.RateDate); err != nil {
			return nil, err
		}
		result[c.Code] = c
	}

	return result, nil
}

func (r *CurrencyRepoPostgres) Create(code string, rate float64, date time.Time) error {
	code = strings.ToUpper(code)

	_, err := r.db.Exec(
		`INSERT INTO currencies (code, rate) VALUES ($1, $2, $3)`,
		code,
		rate,
		date,
	)

	if err != nil {
		return domain.ErrAlreadyExists
	}

	return nil
}

func (r *CurrencyRepoPostgres) UpdateOne(code string, rate float64, date time.Time) error {
	code = strings.ToUpper(code)

	res, err := r.db.Exec(
		`UPDATE currencies SET rate = $1, rate_date = $2 WHERE code = $3`,
		rate, date, code,
	)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *CurrencyRepoPostgres) UpdateAll() error {
	_, err := r.db.Exec(`UPDATE currencies SET rate_date = $1`, time.Now())
	return err
}

func (r *CurrencyRepoPostgres) DeleteAll() error {
	_, err := r.db.Exec(`DELETE FROM currencies`)
	return err
}

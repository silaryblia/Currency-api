package repository

import (
	"Currency-apiNew2/internal/currency/domain"
	"database/sql"
	"strings"

	"math/rand"

	"go.uber.org/zap"
)

type CurrencyRepoPostgres struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewCurrencyRepoPostgres(db *sql.DB, logger *zap.Logger) *CurrencyRepoPostgres {
	return &CurrencyRepoPostgres{db: db, logger: logger}
}

func (r *CurrencyRepoPostgres) GetOne(code string) (float64, error) {
	code = strings.ToUpper(code)

	var rate float64

	err := r.db.QueryRow(
		`SELECT rate FROM currencies WHERE code = $1`,
		code,
	).Scan(&rate)

	if err == sql.ErrNoRows {
		return 0, domain.ErrNotFound
	}

	return rate, err
}

func (r *CurrencyRepoPostgres) GetAll() (map[string]float64, error) {
	rows, err := r.db.Query(`SELECT code, rate FROM currencies`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)

	for rows.Next() {
		var code string
		var rate float64
		if err := rows.Scan(&code, &rate); err != nil {
			return nil, err
		}
		result[code] = rate
	}

	return result, nil
}

func (r *CurrencyRepoPostgres) Create(code string, rate float64) error {
	code = strings.ToUpper(code)

	_, err := r.db.Exec(
		`INSERT INTO currencies (code, rate) VALUES ($1, $2)`,
		code,
		rate,
	)

	if err != nil {
		return domain.ErrAlreadyExists
	}

	return nil
}

func (r *CurrencyRepoPostgres) UpdateOne(code string, rate float64) error {
	code = strings.ToUpper(code)

	res, err := r.db.Exec(
		`UPDATE currencies SET rate = $1 WHERE code = $2`,
		rate,
		code,
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
	rows, err := r.db.Query(`SELECT code, rate FROM currencies`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var code string
		var rate float64
		if err := rows.Scan(&code, &rate); err != nil {
			return err
		}

		change := rand.Float64()*10 - 5
		newVal := rate + change
		if newVal < 0 {
			newVal = 0
		}

		if _, err := r.db.Exec(`UPDATE currencies SET rate = $1 WHERE code = $2`, newVal, code); err != nil {
			return err
		}
	}

	return nil
}

func (r *CurrencyRepoPostgres) DeleteAll() error {
	_, err := r.db.Exec(`DELETE FROM currencies`)
	return err
}

package db

import (
	"database/sql"
)

type Store interface {
	GetDeductionByType(deductionType string) (*Deduction, error)
}

type SQLStore struct {
	db *sql.DB
}

func NewStore(conn *sql.DB) Store {
	return &SQLStore{
		db: conn,
	}
}

func (s *SQLStore) GetDeductionByType(deductionType string) (*Deduction, error) {
	var d Deduction
	stmt, err := s.db.Prepare("SELECT * FROM deductions WHERE type = $1")
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(deductionType).Scan(&d.Type, &d.Amount)
	if err != nil {
		return nil, err
	}

	return &d, err
}

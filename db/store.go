package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	GetAllDeductions(ctx context.Context) ([]Deduction, error)
	GetDeductionByType(ctx context.Context, deductionType string) (*Deduction, error)
}

type SQLStore struct {
	db *sql.DB
}

func NewStore(conn *sql.DB) Store {
	return &SQLStore{
		db: conn,
	}
}

func (c *SQLStore) GetAllDeductions(ctx context.Context) ([]Deduction, error) {
	var deductions []Deduction
	rows, err := c.db.QueryContext(ctx, "SELECT type, amount FROM deductions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d Deduction
		err := rows.Scan(&d.Type, &d.Amount)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			return nil, err
		}

		deductions = append(deductions, d)
	}

	return deductions, nil
}

func (s *SQLStore) GetDeductionByType(ctx context.Context, deductionType string) (*Deduction, error) {
	var d Deduction
	stmt, err := s.db.Prepare("SELECT type, amount FROM deductions WHERE type = $1")
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRowContext(ctx, deductionType).Scan(&d.Type, &d.Amount)
	if err != nil {
		return nil, err
	}

	return &d, err
}

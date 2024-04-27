package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	GetAllDeductions(ctx context.Context) ([]Deduction, error)
	UpdateDeductionByType(ctx context.Context, deductionType string, arg UpdateDeductionParams) (*Deduction, error)
}

type SQLStore struct {
	db *sql.DB
}

func NewStore(conn *sql.DB) Store {
	return &SQLStore{
		db: conn,
	}
}

func (s *SQLStore) GetAllDeductions(ctx context.Context) ([]Deduction, error) {
	var deductions []Deduction
	rows, err := s.db.QueryContext(ctx, "SELECT type, amount FROM deductions")
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

func (s *SQLStore) UpdateDeductionByType(ctx context.Context, deductionType string, arg UpdateDeductionParams) (*Deduction, error) {
	stmt, err := s.db.Prepare(`
		UPDATE deductions
		SET
		    amount = COALESCE($1, amount),
			updated_at = NOW()
		WHERE 
			type = $2
		RETURNING type, amount
	`)
	if err != nil {
		return nil, err
	}

	var d Deduction
	err = stmt.QueryRowContext(ctx, arg.Amount, deductionType).Scan(&d.Type, &d.Amount)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

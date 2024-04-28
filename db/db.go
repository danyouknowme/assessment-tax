package db

import (
	"database/sql"
)

func PrepareDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'deduction_type') THEN
				CREATE TYPE deduction_type AS ENUM ('personal', 'donation', 'k-receipt');
			END IF;
		END $$;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "deductions" (
			"id" SERIAL PRIMARY KEY,
			"type" deduction_type NOT NULL,
			"amount" DECIMAL(10, 2) NOT NULL,
			"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			"updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_deduction_type UNIQUE ("type")
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO "deductions" ("type", "amount")
		VALUES
			('personal', 60000.00),
			('donation', 100000.00),
			('k-receipt', 50000.00)
		ON CONFLICT (type) DO NOTHING
	`)
	if err != nil {
		return err
	}

	return nil
}

func ResetDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		TRUNCATE TABLE "deductions" RESTART IDENTITY CASCADE
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO "deductions" ("type", "amount")
		VALUES
			('personal', 60000.00),
			('donation', 100000.00),
			('k-receipt', 50000.00)
		ON CONFLICT (type) DO NOTHING
	`)
	if err != nil {
		return err
	}

	return nil
}

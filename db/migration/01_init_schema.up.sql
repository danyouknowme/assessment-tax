-- Defined Type
CREATE TYPE IF NOT EXISTS deduction_type AS ENUM ('personal', 'donation', 'k-receipt');

-- Table Definition
CREATE TABLE IF NOT EXISTS "deductions" (
    "id" SERIAL PRIMARY KEY,
    "type" deduction_type NOT NULL,
    "amount" DECIMAL(10, 2) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_deduction_type UNIQUE ("type")
    );

-- Insertions with ON CONFLICT DO NOTHING to avoid duplicates
INSERT INTO "deductions" ("type", "amount")
VALUES
    ('personal', 60000.00),
    ('donation', 100000.00),
    ('k-receipt', 50000.00)
ON CONFLICT (type) DO NOTHING;
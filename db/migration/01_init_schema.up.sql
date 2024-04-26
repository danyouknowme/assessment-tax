-- Defined Type
CREATE TYPE deduction_type AS ENUM ('personal', 'donation', 'k-receipt');

-- Table Definition
CREATE TABLE IF NOT EXISTS "deductions" (
    "id" SERIAL PRIMARY KEY,
    "type" deduction_type NOT NULL,
    "amount" DECIMAL(10, 2) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO "deductions" ("type", "amount") VALUES ('personal', 60000.00);
INSERT INTO "deductions" ("type", "amount") VALUES ('donation', 100000.00);
INSERT INTO "deductions" ("type", "amount") VALUES ('k-receipt', 50000.00);

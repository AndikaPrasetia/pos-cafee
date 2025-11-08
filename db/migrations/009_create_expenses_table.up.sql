-- Expenses table
CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE,
    description VARCHAR(200) NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    category VARCHAR(50) NOT NULL CHECK (category IN ('bahan', 'operasional', 'lainnya')),
    expense_date DATE NOT NULL,
    receipt_number VARCHAR(100),
    notes TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

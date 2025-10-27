CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL,
    tg_id BIGINT UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role TEXT CHECK (role IN ('admin', 'skladchi')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) UNIQUE NOT NULL,
	quantity NUMERIC NOT NULL,
	purchase_price NUMERIC NOT NULL,
    sell_price NUMERIC NOT NULL CHECK (sell_price >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);


CREATE TABLE IF NOT EXISTS stock_levels (
    product_id INT PRIMARY KEY REFERENCES products(id),
    qty NUMERIC NOT NULL DEFAULT 0,
    cost_price NUMERIC NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    qty NUMERIC NOT NULL,
    cost_price NUMERIC NOT NULL,
    total_cost NUMERIC GENERATED ALWAYS AS (qty * cost_price) STORED,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    qty NUMERIC NOT NULL,
    sell_price NUMERIC NOT NULL,
    total_revenue NUMERIC GENERATED ALWAYS AS (qty * sell_price) STORED,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now()
);

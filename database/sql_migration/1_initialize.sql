-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	full_name VARCHAR NOT NULL,
	email VARCHAR NOT NULL,
	password_hash VARCHAR NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT (NOW()),
	updated_at TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE product_categories (
	id SERIAL PRIMARY KEY,
	name VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT (NOW()),
	updated_at TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE products (
	id SERIAL PRIMARY KEY,
	cat_id int,
	name VARCHAR NOT NULL,
	description VARCHAR NOT NULL,
	is_available BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP NOT NULL DEFAULT (NOW()),
	updated_at TIMESTAMP NOT NULL DEFAULT (NOW()),
	CONSTRAINT fk_cat_id 
	FOREIGN KEY (cat_id) REFERENCES product_categories(id) ON DELETE CASCADE
);

CREATE TABLE product_stocks (
	id SERIAL PRIMARY KEY,
	product_id int,
	stock_unit INT NOT NULL,
	price_per_unit INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT (NOW()),
	updated_at TIMESTAMP NOT NULL DEFAULT (NOW()),
	CONSTRAINT fk_product_id 
	FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

DROP TYPE IF EXISTS status;
CREATE TYPE status AS ENUM ('Failed', 'Unpaid', 'Paid', 'Canceled');

CREATE TABLE transactions (
	id SERIAL PRIMARY KEY,
	user_id int,
	product_id int,
	unit int NOT NULL,
	total_price int NOT NULL,
	status status NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT (NOW()),
	updated_at TIMESTAMP NOT NULL DEFAULT (NOW()),
	stock_retreived BOOLEAN DEFAULT FALSE,
	expired_at TIMESTAMPTZ,
	CONSTRAINT fk_user_id 
	FOREIGN KEY (user_id) REFERENCES users(id),
	CONSTRAINT fk_product_id 
	FOREIGN KEY (product_id) REFERENCES products(id)
);

-- +migrate StatementEnd

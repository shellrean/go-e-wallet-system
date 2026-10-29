-- PostgreSQL Schema for Go E-Wallet System

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(150) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    email VARCHAR(150) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email_verified_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    account_number VARCHAR(30) NOT NULL UNIQUE,
    balance NUMERIC(19, 2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS factors (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    pin VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    sof_number VARCHAR(30) NOT NULL,
    dof_number VARCHAR(30) NOT NULL,
    transaction_type CHAR(1) NOT NULL,
    amount NUMERIC(19, 2) NOT NULL,
    transaction_datetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);

CREATE TABLE IF NOT EXISTS topup (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status SMALLINT NOT NULL DEFAULT 0,
    amount NUMERIC(19, 2) NOT NULL,
    snap_url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS templates (
    code VARCHAR(64) PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    body TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(150) NOT NULL,
    body TEXT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    is_read SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);

-- Default notification templates
INSERT INTO templates (code, title, body) VALUES
    ('TRANSFER', 'Transfer Berhasil', 'Transfer sejumlah Rp {{.amount}} berhasil dilakukan.'),
    ('TRANSFER_DES', 'Dana Masuk', 'Anda menerima transfer sejumlah Rp {{.amount}}.'),
    ('TOPUP_SUCCESS', 'Top-Up Berhasil', 'Top-Up sejumlah Rp {{.amount}} telah berhasil masuk ke saldo Anda.')
ON CONFLICT (code) DO NOTHING;

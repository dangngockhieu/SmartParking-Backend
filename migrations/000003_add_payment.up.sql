ALTER TABLE users
ADD COLUMN money BIGINT NOT NULL DEFAULT 0;

CREATE TABLE wallet_transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,

    type VARCHAR(20) NOT NULL,
    amount BIGINT NOT NULL,

    balance_before BIGINT NOT NULL DEFAULT 0,
    balance_after BIGINT NOT NULL DEFAULT 0,

    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    order_code BIGINT UNSIGNED UNIQUE,
    transaction_id VARCHAR(100),
    description TEXT,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_wallet_transactions_user
        FOREIGN KEY (user_id) REFERENCES users(id),

    CONSTRAINT chk_wallet_type
        CHECK (type IN ('DEPOSIT', 'DEDUCT')),

    CONSTRAINT chk_wallet_status
        CHECK (status IN ('PENDING', 'SUCCESS', 'FAILED', 'CANCELED')),

    CONSTRAINT chk_wallet_amount_positive
        CHECK (amount > 0),

    CONSTRAINT chk_wallet_balance_non_negative
        CHECK (balance_before >= 0 AND balance_after >= 0)
);

CREATE INDEX idx_wallet_transactions_user_created_id
ON wallet_transactions (user_id, created_at DESC, id DESC);

CREATE INDEX idx_wallet_transactions_created_id
ON wallet_transactions (created_at DESC, id DESC);

CREATE INDEX idx_wallet_transactions_status_created_at
ON wallet_transactions (status, created_at);

UPDATE parking_sessions
SET fee = ROUND(fee);

ALTER TABLE parking_sessions
MODIFY COLUMN fee BIGINT NOT NULL DEFAULT 0;
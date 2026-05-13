ALTER TABLE users
DROP COLUMN IF EXISTS money;

DROP INDEX IF EXISTS idx_wallet_transactions_success_created_at;

DROP INDEX IF EXISTS idx_wallet_transactions_created_id;

DROP INDEX IF EXISTS idx_wallet_transactions_user_created_id;

DROP TABLE IF EXISTS wallet_transactions;

ALTER TABLE parking_sessions
MODIFY COLUMN fee DOUBLE NOT NULL DEFAULT 0;


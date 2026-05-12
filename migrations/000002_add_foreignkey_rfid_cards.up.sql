ALTER TABLE rfid_cards
DROP COLUMN owner_name;

ALTER TABLE rfid_cards
ADD COLUMN user_id BIGINT UNSIGNED NULL;
ALTER TABLE rfid_cards
ADD CONSTRAINT fk_rfid_cards_user_id
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE SET NULL
ON UPDATE CASCADE;
ALTER TABLE rfid_cards
ADD CONSTRAINT uq_rfid_cards_user_id UNIQUE (user_id);

ALTER TABLE slot_histories
DROP FOREIGN KEY fk_slot_histories_slot;

ALTER TABLE slot_histories
DROP FOREIGN KEY fk_slot_histories_user;

DROP INDEX idx_slot_histories_slot_id ON slot_histories;
DROP INDEX idx_slot_histories_user_id ON slot_histories;
DROP INDEX idx_slot_histories_created_at ON slot_histories;
DROP INDEX idx_slot_histories_slot_id_created_at ON slot_histories;

DROP TABLE IF EXISTS slot_histories;
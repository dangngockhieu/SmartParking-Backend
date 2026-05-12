ALTER TABLE rfid_cards
DROP FOREIGN KEY fk_rfid_cards_user_id;

ALTER TABLE rfid_cards
DROP INDEX uq_rfid_cards_user_id;

ALTER TABLE rfid_cards
DROP COLUMN user_id;

ALTER TABLE rfid_cards
ADD COLUMN owner_name VARCHAR(255) NULL;

CREATE TABLE slot_histories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    slot_id BIGINT UNSIGNED NOT NULL,

    old_device VARCHAR(50) NULL,
    new_device VARCHAR(50) NULL,

    old_port INT NULL,
    new_port INT NULL,

    action ENUM(
        'DEVICE_CHANGE',
        'STATUS_CHANGE',
        'SYSTEM_FIX',
        'MAINTAIN_MODE'
    ) NOT NULL DEFAULT 'DEVICE_CHANGE',

    user_id BIGINT UNSIGNED NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_slot_histories_slot
        FOREIGN KEY (slot_id)
        REFERENCES parking_slots(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_slot_histories_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
);

CREATE INDEX idx_slot_histories_slot_id
ON slot_histories(slot_id);

CREATE INDEX idx_slot_histories_user_id
ON slot_histories(user_id);

CREATE INDEX idx_slot_histories_created_at
ON slot_histories(created_at);

CREATE INDEX idx_slot_histories_slot_id_created_at
ON slot_histories(slot_id, created_at);
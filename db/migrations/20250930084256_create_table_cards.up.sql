CREATE TABLE cards (
    id          VARCHAR(100) PRIMARY KEY,
    board_id    VARCHAR(100) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    is_closed   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    CONSTRAINT fk_cards_board FOREIGN KEY (board_id)
        REFERENCES boards (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

-- function untuk auto update kolom updated_at
CREATE OR REPLACE FUNCTION update_cards_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- trigger pasang ke tabel cards
CREATE TRIGGER update_cards_updated_at
BEFORE UPDATE ON cards
FOR EACH ROW
EXECUTE FUNCTION update_cards_updated_at_column();

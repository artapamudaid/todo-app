CREATE TABLE boards (
    id          VARCHAR(100) PRIMARY KEY,
    project_id  VARCHAR(100) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    CONSTRAINT fk_boards_project FOREIGN KEY (project_id)
        REFERENCES projects (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

-- function untuk auto update kolom updated_at
CREATE OR REPLACE FUNCTION update_boards_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- trigger pasang ke tabel boards
CREATE TRIGGER update_boards_updated_at
BEFORE UPDATE ON boards
FOR EACH ROW
EXECUTE FUNCTION update_boards_updated_at_column();

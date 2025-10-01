CREATE TABLE project_users (
    id          VARCHAR(100) PRIMARY KEY,
    project_id  VARCHAR(100) NOT NULL,
    user_id     VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    CONSTRAINT fk_project_users_project FOREIGN KEY (project_id)
        REFERENCES projects (id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_project_users_user FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

-- function untuk auto update kolom updated_at
CREATE OR REPLACE FUNCTION update_project_users_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- trigger pasang ke tabel project_users
CREATE TRIGGER update_project_users_updated_at
BEFORE UPDATE ON project_users
FOR EACH ROW
EXECUTE FUNCTION update_project_users_updated_at_column();

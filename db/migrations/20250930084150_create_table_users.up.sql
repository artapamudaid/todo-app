CREATE TABLE users (
    id             VARCHAR(100) PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    email          VARCHAR(100) NOT NULL UNIQUE,
    password       VARCHAR(100) NOT NULL,
    token          VARCHAR(500) NULL,
    role_id        VARCHAR(100) NOT NULL,
    department_id  VARCHAR(100) NOT NULL,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP NULL,
    CONSTRAINT fk_users_role FOREIGN KEY (role_id)
        REFERENCES roles (id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_users_department FOREIGN KEY (department_id)
        REFERENCES departments (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

-- function untuk auto update kolom updated_at
CREATE OR REPLACE FUNCTION update_users_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- trigger pasang ke tabel users
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_updated_at_column();

ALTER TABLE projects
ADD COLUMN IF NOT EXISTS department_id VARCHAR(100),
ADD CONSTRAINT fk_projects_department
    FOREIGN KEY (department_id) REFERENCES departments(id)
    ON UPDATE CASCADE ON DELETE SET NULL;
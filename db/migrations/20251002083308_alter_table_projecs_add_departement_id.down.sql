ALTER TABLE projects
		DROP CONSTRAINT IF EXISTS fk_projects_department,
		DROP COLUMN IF EXISTS department_id;
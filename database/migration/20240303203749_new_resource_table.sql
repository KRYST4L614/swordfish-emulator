-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resource (
 id TEXT PRIMARY KEY,
 data JSONB
);

CREATE TABLE IF NOT EXISTS operation_privilege (
 resource TEXT PRIMARY KEY,
 data JSONB
);

CREATE TABLE IF NOT EXISTS privilege (
 id SERIAL PRIMARY KEY,
 name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS role (
 id SERIAL PRIMARY KEY,
 name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS role_privilege (
 role_id INT REFERENCES role(id) ON DELETE CASCADE,
 privilege_id INT REFERENCES privilege(id) ON DELETE CASCADE,
 PRIMARY KEY (role_id, privilege_id)
);

CREATE TABLE IF NOT EXISTS "user" (
 name TEXT PRIMARY KEY,
 password TEXT NOT NULL,
 role_id INT REFERENCES role(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resource;
DROP TABLE IF EXISTS operation_privilege;
DROP TABLE IF EXISTS privilege;
DROP TABLE IF EXISTS role;
DROP TABLE IF EXISTS role_privilege;
DROP TABLE IF EXISTS user;
-- +goose StatementEnd

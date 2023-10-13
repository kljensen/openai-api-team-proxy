-- Schema for sqlite3 database

-- Table: users
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255)
);

-- Table: roles
CREATE TABLE roles (
    name TEXT PRIMARY KEY
        CHECK (name REGEXP '^[a-zA-Z0-9_-]+$') AND length(name) <= 100
        UNIQUE
);

-- Table: user_roles
CREATE TABLE user_roles (
    user_id INTEGER NOT NULL,
    role_name TEXT NOT NULL,
    PRIMARY KEY (user_id, role_name),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (role_name) REFERENCES roles (name)
);


-- Table: permissions
CREATE TABLE permissions (
    name TEXT PRIMARY KEY
        CHECK (name REGEXP '^[a-zA-Z0-9_-]+$') AND length(name) <= 100
        UNIQUE,
    description TEXT,
    path TEXT NOT NULL,
    rule JSON
);

-- Check that rule is valid JSON
CREATE TRIGGER permissions_rule_check
BEFORE INSERT ON permissions
FOR EACH ROW
BEGIN
    SELECT CASE
        WHEN (NEW.rule IS NULL OR json_valid(NEW.rule)) THEN
            NULL
        ELSE
            RAISE(ABORT, 'Invalid JSON rule')
        END;
END;

-- Table: role_permissions
CREATE TABLE role_permissions (
    role_name TEXT NOT NULL,
    permission_name TEXT NOT NULL,
    PRIMARY KEY (role_name, permission_name),
    FOREIGN KEY (role_name) REFERENCES roles (name),
    FOREIGN KEY (permission_name) REFERENCES permissions (name)
);

-- Insert some common rules

-- Can use only gpt-3.5-turbo or gpt-3.5-turbo-instruct on /v1/completions
INSERT INTO permissions (name, description, path, rule) VALUES (
    'completions-only-gpt-3.5',
    'Can use only gpt-3.5-turbo or gpt-3.5-turbo-instruct on for chat completion',
    '/v1/completions',
    -- Use JSONLogic. 
    -- See https://jsonlogic.com/operations.html
    '{
        "or": [
            {"==": [{"var": "model"}, "gpt-3.5-turbo"]},
            {"==": [{"var": "model"}, "gpt-3.5-turbo-instruct"]}
        ]
    }'
);
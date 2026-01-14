
-- +goose Up
CREATE TABLE IF NOT EXISTS employees (
    id            SERIAL          PRIMARY KEY,
    user_id       BIGINT          UNIQUE NOT NULL,      -- each user → at most one employee

    job_title     VARCHAR(100)    NOT NULL,
    country       VARCHAR(100)    NOT NULL,
    salary        DECIMAL(12,2)   NOT NULL,
    created_at    TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,

    -- Constraint for "Foreign Key"
    CONSTRAINT fk_employee_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS employees;
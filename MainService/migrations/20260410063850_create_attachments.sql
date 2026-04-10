-- +goose Up
CREATE TABLE attachments (
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(255) NOT NULL,
                             url VARCHAR(255) NOT NULL,
                             lesson_id INTEGER NOT NULL,
                             created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS attachments;

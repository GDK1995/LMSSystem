-- +goose Up
CREATE TABLE user_lessons (
                              user_id TEXT NOT NULL,
                              lesson_id INTEGER NOT NULL,
                              PRIMARY KEY (user_id, lesson_id)
);

-- +goose Down
DROP TABLE IF EXISTS user_lessons;
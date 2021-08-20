-- +goose Up

CREATE TABLE IF NOT EXISTS "questions" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" TEXT,
    "options" TEXT,
    "correct_option" INTEGER,
    "quiz" INTEGER NOT NULL,
    "points" INTEGER,
    FOREIGN KEY ("quiz") REFERENCES "quiz"("id")

);

-- +goose Down

DROP TABLE IF EXISTS questions;
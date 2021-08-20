-- +goose Up

CREATE TABLE IF NOT EXISTS "quiz" (
    "id" INTEGER  PRIMARY KEY AUTOINCREMENT,
    "name" TEXT,
    "description" TEXT
);

-- +goose Down

DROP TABLE IF EXISTS quiz;
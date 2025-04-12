CREATE TABLE IF NOT EXISTS "roles" (
    "id" serial PRIMARY KEY,
    "title" varchar NOT NULL UNIQUE
);
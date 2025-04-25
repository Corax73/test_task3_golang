CREATE TABLE IF NOT EXISTS "roles" (
    "id" serial PRIMARY KEY,
    "title" varchar NOT NULL UNIQUE,
    "abilities" JSONB,
    "created_at" date NOT NULL DEFAULT 'now()'
);
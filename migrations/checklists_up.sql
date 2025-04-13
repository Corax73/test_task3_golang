CREATE TABLE IF NOT EXISTS "checklists" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL,
    "name" varchar NOT NULL DEFAULT '',
    "created_at" date NOT NULL DEFAULT 'now()',
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
CREATE INDEX ON "checklists" ("user_id");
CREATE TABLE IF NOT EXISTS "users" (
    "id" serial PRIMARY KEY,
    "role_id" int NOT NULL,
    "login" varchar NOT NULL DEFAULT '',
    "email" varchar NOT NULL DEFAULT '',
    "password" varchar NOT NULL DEFAULT '',
    "created_at" date NOT NULL DEFAULT 'now()',
    FOREIGN KEY ("role_id") REFERENCES "roles" ("id")
);
CREATE INDEX ON "users" ("role_id");
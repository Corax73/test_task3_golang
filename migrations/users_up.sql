CREATE TABLE IF NOT EXISTS "users" (
    "id" serial PRIMARY KEY,
    "role_id" int NOT NULL,
    "login" varchar NOT NULL DEFAULT '' UNIQUE,
    "email" varchar NOT NULL DEFAULT '' UNIQUE,
    "password" varchar NOT NULL DEFAULT '',
    "checklists_quantity" int DEFAULT 0 NOT NULL,
    "created_at" date NOT NULL DEFAULT 'now()',
    FOREIGN KEY ("role_id") REFERENCES "roles" ("id")
);
CREATE INDEX ON "users" ("role_id");